package dues

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/resp"
	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/timediff"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
)

var (
	ErrDuesNotFound  = errors.New("tagihan iuran tidak ditemukan")
	ErrDateInThePast = errors.New("tanggal tidak boleh di masa lalu")
	ErrDateFormat    = errors.New("format tanggal tidak sesuai <tahun>-<bulan>-<hari>")
	ErrExistDues     = errors.New("tagihan iuran untuk bulan ini sudah ada")
	ErrProcessedDues = errors.New("seseorang telah melalukan pembayaran untuk tagihan iuran ini, tidak dapat dimodifikasi")
)

type (
	AddDuesIn struct {
		Date      string `json:"date"`
		IdrAmount string `json:"idr_amount"`
	}
	AddDuesRes struct {
		Id int64 `json:"id"`
	}
	AddDuesOut struct {
		resp.Response
		Res AddDuesRes
	}
)

func (d *DuesDeps) AddDues(ctx context.Context, in AddDuesIn) (out AddDuesOut) {
	var err error
	out.Response = resp.NewResponse(http.StatusCreated, "", nil)

	if err = ValidateAddDuesIn(in); err != nil {
		out.Response = resp.NewResponse(http.StatusUnprocessableEntity, "", err)
		return
	}

	date, err := time.Parse("2006-01-02", in.Date)
	if err != nil {
		out.Response = resp.NewResponse(http.StatusUnprocessableEntity, "", ErrDateFormat)
		return
	}

	if timediff.IsPast(date) {
		out.Response = resp.NewResponse(http.StatusBadRequest, "", ErrDateInThePast)
		return
	}

	dues, err := d.DuesRepository.FindOtherByYYYYMM(ctx, 0, date)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "find dues by yyyy mm"))
		return
	}

	if dues.Id != 0 {
		out.Response = resp.NewResponse(http.StatusBadRequest, "", ErrExistDues)
		return
	}

	dues = DuesModel{
		Date:      date,
		IdrAmount: in.IdrAmount,
	}

	if dues, err = d.DuesRepository.Save(ctx, dues); err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "save dues"))
		return
	}

	if err = d.MemberDuesRepository.GenerateDues(ctx, dues.Id); err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "generate dues"))
		return
	}

	out.Res.Id = int64(dues.Id)

	return
}

type (
	DuesOut struct {
		Id        int64  `json:"id"`
		Date      string `json:"date"`
		IdrAmount string `json:"idr_amount"`
	}
	QueryDuesRes struct {
		Cursor int64     `json:"cursor"`
		Dues   []DuesOut `json:"dues"`
	}
	QueryDuesOut struct {
		resp.Response
		Res QueryDuesRes
	}
)

func (d *DuesDeps) QueryDues(ctx context.Context, cursor, limit string) (out QueryDuesOut) {
	var err error
	out.Response = resp.NewResponse(http.StatusOK, "", nil)

	fromCursor, _ := strconv.ParseInt(cursor, 10, 64)
	nlimit, _ := strconv.ParseInt(limit, 10, 64)
	if nlimit == 0 {
		nlimit = 25
	}

	dues, err := d.DuesRepository.Query(ctx, fromCursor, nlimit)
	if err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "query dues"))
		return
	}

	duesLen := len(dues)

	var nextCursor int64
	if duesLen != 0 {
		nextCursor = int64(dues[duesLen-1].Id)
	}

	outDues := make([]DuesOut, duesLen)
	for i, d := range dues {
		outDues[i] = DuesOut{
			Id:        int64(d.Id),
			Date:      d.Date.Format("2006-01"),
			IdrAmount: d.IdrAmount,
		}
	}

	out.Res = QueryDuesRes{
		Cursor: nextCursor,
		Dues:   outDues,
	}

	return
}

type (
	EditDuesIn struct {
		Date      string `json:"date"`
		IdrAmount string `json:"idr_amount"`
	}
	EditDuesRes struct {
		Id int64 `json:"id"`
	}
	EditDuesOut struct {
		resp.Response
		Res EditDuesRes
	}
)

func (d *DuesDeps) EditDues(ctx context.Context, pid string, in EditDuesIn) (out EditDuesOut) {
	var err error
	out.Response = resp.NewResponse(http.StatusOK, "", nil)

	id, err := strconv.ParseUint(pid, 10, 64)
	if err != nil {
		out.Response = resp.NewResponse(http.StatusNotFound, "", ErrDuesNotFound)
		return
	}

	if err = ValidateEditDuesIn(in); err != nil {
		out.Response = resp.NewResponse(http.StatusUnprocessableEntity, "", errors.Wrap(err, "edit dues validation"))
		return
	}

	date, err := time.Parse("2006-01-02", in.Date)
	if err != nil {
		out.Response = resp.NewResponse(http.StatusUnprocessableEntity, "", ErrDateFormat)
		return
	}

	if timediff.IsPast(date) {
		out.Response = resp.NewResponse(http.StatusBadRequest, "", ErrDateInThePast)
		return
	}

	dues, err := d.DuesRepository.FindById(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		out.Response = resp.NewResponse(http.StatusNotFound, "", ErrDuesNotFound)
		return
	}
	if err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "find document by id"))
		return
	}

	otherDues, err := d.DuesRepository.FindOtherByYYYYMM(ctx, dues.Id, date)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "find dues by yyyy mm"))
		return
	}

	if otherDues.Id != 0 {
		out.Response = resp.NewResponse(http.StatusBadRequest, "", ErrExistDues)
		return
	}

	duesList, err := d.MemberDuesRepository.CheckSomeonePaid(ctx, dues.Id)
	if err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "check someone paid dues"))
		return
	}

	if len(duesList) != 0 {
		out.Response = resp.NewResponse(http.StatusBadRequest, "", ErrProcessedDues)
		return
	}

	dues.Date = date
	dues.IdrAmount = in.IdrAmount

	if err = d.DuesRepository.UpdateById(ctx, id, dues); err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "update dues by id"))
		return
	}

	out.Res.Id = int64(id)

	return
}

type (
	RemoveDuesRes struct {
		Id int64 `json:"id"`
	}
	RemoveDuesOut struct {
		resp.Response
		Res RemoveDuesRes
	}
)

func (d *DuesDeps) RemoveDues(ctx context.Context, pid string) (out RemoveDuesOut) {
	var err error
	out.Response = resp.NewResponse(http.StatusOK, "", nil)

	id, err := strconv.ParseUint(pid, 10, 64)
	if err != nil {
		out.Response = resp.NewResponse(http.StatusNotFound, "", ErrDuesNotFound)
		return
	}

	dues, err := d.DuesRepository.FindById(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		out.Response = resp.NewResponse(http.StatusNotFound, "", ErrDuesNotFound)
		return
	}
	if err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "find dues by id"))
		return
	}

	duesList, err := d.MemberDuesRepository.CheckSomeonePaid(ctx, dues.Id)
	if err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "check someone paid dues"))
		return
	}

	if len(duesList) != 0 {
		out.Response = resp.NewResponse(http.StatusBadRequest, "", ErrProcessedDues)
		return
	}

	if err = d.DuesRepository.DeleteById(ctx, id); err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "delete dues by id"))
		return
	}

	if err = d.MemberDuesRepository.DeleteByDuesId(ctx, id); err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "delete member dues by dues id"))
		return
	}

	out.Res.Id = int64(id)

	return
}

type (
	CheckPaidDuesRes struct {
		IsPaid bool `json:"is_paid"`
	}
	CheckPaidDuesOut struct {
		resp.Response
		Res CheckPaidDuesRes
	}
)

func (d *DuesDeps) CheckPaidDues(ctx context.Context, pid string) (out CheckPaidDuesOut) {
	out.StatusCode = http.StatusOK

	id, err := strconv.ParseUint(pid, 10, 64)
	if err != nil {
		out.Response = resp.NewResponse(http.StatusNotFound, "", ErrDuesNotFound)
		return
	}

	dues, err := d.DuesRepository.FindById(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		out.Response = resp.NewResponse(http.StatusNotFound, "", ErrDuesNotFound)
		return
	}
	if err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "find dues by id"))
		return
	}

	duesList, err := d.MemberDuesRepository.CheckSomeonePaid(ctx, dues.Id)
	if err != nil {
		out.Response = resp.NewResponse(http.StatusBadRequest, "", ErrProcessedDues)
		return
	}

	isPaid := false
	if len(duesList) != 0 {
		isPaid = true
	}

	out.Res.IsPaid = isPaid

	return
}
