package dues

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/cashflow"
	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/httpdecode"
	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/resp"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"gopkg.in/guregu/null.v4"
)

type (
	MemberDuesOut struct {
		Id           int64  `json:"id"`
		DuesId       int64  `json:"dues_id"`
		Date         string `json:"date"`
		Status       string `json:"status"`
		IdrAmout     string `json:"idr_amount"`
		ProveFileUrl string `json:"prove_file_url"`
	}
	MemberDuesRes struct {
		TotalDues  string          `json:"total_dues"`
		PaidDues   string          `json:"paid_dues"`
		UnpaidDues string          `json:"unpaid_dues"`
		Cursor     int64           `json:"cursor"`
		Dues       []MemberDuesOut `json:"dues"`
	}
	QueryMemberDuesOut struct {
		resp.Response
		Res MemberDuesRes
	}
)

func (d *DuesDeps) QueryMemberDues(ctx context.Context, uid string, cursor string) (out QueryMemberDuesOut) {
	var err error
	out.Response = resp.NewResponse(http.StatusOK, "", nil)

	_, err = uuid.FromString(uid)
	if err != nil {
		out.Response = resp.NewResponse(http.StatusBadRequest, "", errors.Wrap(err, "member id not valid"))
		return
	}

	_, err = d.MemberRepository.FindById(ctx, uid)
	if errors.Is(err, pgx.ErrNoRows) {
		out.Response = resp.NewResponse(http.StatusNotFound, "", errors.Wrap(err, "no row find member by id"))
		return
	}

	if err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "find member by id"))
		return
	}

	duefFlow := func(ctx context.Context, uid string, status DuesStatus, dues chan float64, res chan resp.Response) {
		var r resp.Response
		amts, err := d.DuesRepository.QueryAmtByUidStatus(ctx, uid, status)
		if err != nil {
			r = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "query "+status.String+" dues by uid"))
		}

		var namts float64
		for _, u := range amts {
			cash, _ := strconv.ParseFloat(u, 64)
			namts += cash
		}

		dues <- namts
		res <- r
	}

	paidDues := make(chan float64)
	pRes := make(chan resp.Response)
	go duefFlow(ctx, uid, Paid, paidDues, pRes)

	unpaidDues := make(chan float64)
	uRes := make(chan resp.Response)
	go duefFlow(ctx, uid, Unpaid, unpaidDues, uRes)

	outMemberDues := make(chan []MemberDuesOut)
	nextCursor := make(chan int64)
	qRes := make(chan resp.Response)

	go func(ctx context.Context, uid, cursor string, md chan []MemberDuesOut, nc chan int64, res chan resp.Response) {
		var r resp.Response
		fromCursor, _ := strconv.ParseInt(cursor, 10, 64)
		memberdues, err := d.MemberDuesRepository.QueryMDVByUid(ctx, uid, fromCursor, 25)
		if err != nil {
			r = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "query member dues view by uid"))
		}

		mdsLen := len(memberdues)

		var nextCursor int64
		if mdsLen != 0 {
			nextCursor = int64(memberdues[mdsLen-1].Id)
		}

		outMemberDues := make([]MemberDuesOut, mdsLen)
		for i, d := range memberdues {
			status := d.Status
			if status == Unknown {
				continue
			}

			outMemberDues[i] = MemberDuesOut{
				Id:           int64(d.Id),
				DuesId:       int64(d.DuesId),
				Date:         d.Date.Format("2006-01-02"),
				Status:       status.String,
				IdrAmout:     d.IdrAmount,
				ProveFileUrl: d.ProveFileUrl,
			}
		}

		md <- outMemberDues
		nc <- nextCursor
		res <- r
	}(ctx, uid, cursor, outMemberDues, nextCursor, qRes)

	paidDuesV := <-paidDues
	pRV := <-pRes
	unpaidDuesV := <-unpaidDues
	uRV := <-uRes
	outMemberDuesV := <-outMemberDues
	nextCursorV := <-nextCursor
	qRV := <-qRes

	if pRV.Error != nil {
		out.Response = pRV
		return
	}

	if uRV.Error != nil {
		out.Response = uRV
		return
	}

	if qRV.Error != nil {
		out.Response = qRV
		return
	}

	out.Res = MemberDuesRes{
		TotalDues:  strconv.FormatFloat(paidDuesV, 'f', -1, 64),
		PaidDues:   strconv.FormatFloat(paidDuesV, 'f', -1, 64),
		UnpaidDues: strconv.FormatFloat(unpaidDuesV, 'f', -1, 64),
		Cursor:     nextCursorV,
		Dues:       outMemberDuesV,
	}

	return
}

type (
	MembersDuesOut struct {
		Id            int64  `json:"id"`
		MemberId      string `json:"member_id"`
		Status        string `json:"status"`
		Name          string `json:"name"`
		ProfilePicUrl string `json:"profile_pic_url"`
	}
	QueryMembersDuesRes struct {
		Cursor     int64            `json:"cursor"`
		MemberDues []MembersDuesOut `json:"member_dues"`
	}
	QueryMembersDuesOut struct {
		resp.Response
		Res QueryMembersDuesRes
	}
)

func (d *DuesDeps) QueryMembersDues(ctx context.Context, pid, cursor string) (out QueryMembersDuesOut) {
	var err error
	out.Response = resp.NewResponse(http.StatusOK, "", nil)

	id, err := strconv.ParseUint(pid, 10, 64)
	if err != nil {
		out.StatusCode = http.StatusBadRequest
		err = errors.Wrap(err, "parse uint")
		return
	}

	var dues DuesModel
	if id == 0 {
		dues, err = d.DuesRepository.Latest(ctx)
		if err != nil {
			out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "latest dues"))
			return
		}
	}

	if id != 0 {
		dues, err = d.DuesRepository.FindById(ctx, id)
		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "find dues by id"))
			return
		}
	}

	if dues.Id == 0 {
		out.Res = QueryMembersDuesRes{
			Cursor:     0,
			MemberDues: make([]MembersDuesOut, 0, 0),
		}
		return
	}

	fromCursor, _ := strconv.ParseInt(cursor, 10, 64)
	memberDues, err := d.MemberDuesRepository.QueryDMVByDuesId(ctx, dues.Id, fromCursor, 25)
	if err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "query member dues by uid"))
		return
	}

	mdsLen := len(memberDues)

	var nextCursor int64
	if mdsLen != 0 {
		nextCursor = int64(memberDues[mdsLen-1].Id)
	}

	outMembersDues := make([]MembersDuesOut, mdsLen)
	for i, m := range memberDues {
		status := m.Status
		if status == Unknown {
			continue
		}

		outMembersDues[i] = MembersDuesOut{
			Id:            int64(m.Id),
			MemberId:      m.MemberId,
			Status:        status.String,
			Name:          m.Name,
			ProfilePicUrl: m.ProfilePicUrl,
		}
	}

	out.Res = QueryMembersDuesRes{
		Cursor:     nextCursor,
		MemberDues: outMembersDues,
	}

	return
}

type (
	PayMemberDuesIn struct {
		File httpdecode.FileHeader `mapstructure:"file"`
	}
	PayMemberDuesRes struct {
		Id int64 `json:"id"`
	}
	PayMemberDuesOut struct {
		resp.Response
		Res PayMemberDuesRes
	}
)

func (d *DuesDeps) PayMemberDues(ctx context.Context, uid, pid string, in PayMemberDuesIn) (out PayMemberDuesOut) {
	var err error
	out.Response = resp.NewResponse(http.StatusOK, "", nil)

	_, err = uuid.FromString(uid)
	if err != nil {
		out.Response = resp.NewResponse(http.StatusBadRequest, "", errors.Wrap(err, "member id not valid"))
		return
	}

	id, err := strconv.ParseUint(pid, 10, 64)
	if err != nil {
		out.Response = resp.NewResponse(http.StatusBadRequest, "", errors.Wrap(err, "parse uint"))
		return
	}

	if err = ValidatePayMemberDuesIn(in); err != nil {
		out.Response = resp.NewResponse(http.StatusUnprocessableEntity, "", errors.Wrap(err, "pay member dues validation"))
		return
	}

	memberDues, err := d.MemberDuesRepository.FindUnpaidByIdAndMemberId(ctx, id, uid)
	if errors.Is(err, pgx.ErrNoRows) {
		out.Response = resp.NewResponse(http.StatusNotFound, "", errors.Wrap(err, "no row find dues by id"))
		return
	}
	if err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "find dues by id"))
		return
	}

	var file httpdecode.File
	if in.File.File != nil {
		file = in.File.File
	}

	defer func() {
		if file != nil {
			file.Close()
		}
	}()

	var fileUrl string
	if file != nil {
		filename := strconv.FormatInt(time.Now().Unix(), 10) + "-" + strings.Trim(in.File.Filename, " ")
		if fileUrl, err = d.Upload(filename, file); err != nil {
			out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "upload file"))
			return
		}
	}

	memberDues.Status = Waiting
	memberDues.ProveFileUrl = fileUrl

	if err = d.MemberDuesRepository.UpdateById(ctx, id, memberDues); err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "save member dues"))
		return
	}

	out.Res.Id = int64(memberDues.Id)

	return
}

type (
	EditMemberDuesIn struct {
		File httpdecode.FileHeader `mapstructure:"file"`
	}
	EditMemberDuesRes struct {
		Id int64 `json:"id"`
	}
	EditMemberDuesOut struct {
		resp.Response
		Res EditMemberDuesRes
	}
)

func (d *DuesDeps) EditMemberDues(ctx context.Context, pid string, in EditMemberDuesIn) (out EditMemberDuesOut) {
	var err error
	out.Response = resp.NewResponse(http.StatusOK, "", nil)

	id, err := strconv.ParseUint(pid, 10, 64)
	if err != nil {
		out.Response = resp.NewResponse(http.StatusBadRequest, "", errors.Wrap(err, "parse uint"))
		return
	}

	if err = ValidateEditMemberDuesIn(in); err != nil {
		out.Response = resp.NewResponse(http.StatusUnprocessableEntity, "", errors.Wrap(err, "edit member dues validation"))
		return
	}

	memberDues, err := d.MemberDuesRepository.FindUnpaidById(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		out.Response = resp.NewResponse(http.StatusNotFound, "", errors.Wrap(err, "no row find member dues by id"))
		return
	}
	if err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "find document by id"))
		return
	}

	var file httpdecode.File
	if in.File.File != nil {
		file = in.File.File
	}

	defer func() {
		if file != nil {
			file.Close()
		}
	}()

	var fileUrl string
	if file != nil {
		filename := strconv.FormatInt(time.Now().Unix(), 10) + "-" + strings.Trim(in.File.Filename, " ")
		if fileUrl, err = d.Upload(filename, file); err != nil {
			out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "upload file"))
			return
		}
	}

	memberDues.ProveFileUrl = fileUrl

	if err = d.MemberDuesRepository.UpdateById(ctx, id, memberDues); err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "update document by id"))
		return
	}

	out.Res.Id = int64(id)

	return
}

type (
	PaidMemberDuesIn struct {
		IsPaid null.Bool `json:"is_paid"`
	}
	PaidMemberDuesRes struct {
		Id int64 `json:"id"`
	}
	PaidMemberDuesOut struct {
		resp.Response
		Res EditMemberDuesRes
	}
)

func (d *DuesDeps) PaidMemberDues(ctx context.Context, pid string, in PaidMemberDuesIn) (out PaidMemberDuesOut) {
	var err error
	out.Response = resp.NewResponse(http.StatusOK, "", nil)

	id, err := strconv.ParseUint(pid, 10, 64)
	if err != nil {
		out.Response = resp.NewResponse(http.StatusBadRequest, "", errors.Wrap(err, "parse uint"))
		return
	}

	if err = ValidatePaidMemberDuesIn(in); err != nil {
		out.Response = resp.NewResponse(http.StatusUnprocessableEntity, "", errors.Wrap(err, "edit member dues validation"))
		return
	}

	memberDues, err := d.MemberDuesRepository.FindUnpaidById(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		out.Response = resp.NewResponse(http.StatusNotFound, "", errors.Wrap(err, "no row find member dues by id"))
		return
	}
	if err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "find document by id"))
		return
	}

	dues, err := d.DuesRepository.FindById(ctx, memberDues.DuesId)
	if errors.Is(err, pgx.ErrNoRows) {
		out.Response = resp.NewResponse(http.StatusNotFound, "", errors.Wrap(err, "no row find dues by id"))
		return
	}
	if err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "find document by id"))
		return
	}

	member, err := d.MemberRepository.FindById(ctx, memberDues.MemberId)
	if errors.Is(err, pgx.ErrNoRows) {
		out.Response = resp.NewResponse(http.StatusNotFound, "", errors.Wrap(err, "no row find member by id"))
		return
	}
	if err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "find member by id"))
		return
	}

	if in.IsPaid.Valid && in.IsPaid.Bool == true {
		memberDues.Status = Paid
	}

	if err = d.MemberDuesRepository.UpdateById(ctx, id, memberDues); err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "update document by id"))
		return
	}

	cashflow := cashflow.CashflowModel{
		Date:         memberDues.CreatedAt,
		IdrAmount:    dues.IdrAmount,
		Type:         cashflow.Income,
		Note:         "Pembayaran Iuran Anggota, Nama " + member.Name,
		ProveFileUrl: memberDues.ProveFileUrl,
	}

	if cashflow, err = d.CashflowRepository.Save(ctx, cashflow); err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "save cashflow"))
		return
	}

	out.Res.Id = int64(id)

	return
}
