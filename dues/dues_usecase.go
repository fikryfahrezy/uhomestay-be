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
		out.Response = resp.NewResponse(http.StatusUnprocessableEntity, "", errors.Wrap(err, "add dues validation"))
		return
	}

	date, err := time.Parse("2006-01-02", in.Date)
	if err != nil {
		out.Response = resp.NewResponse(http.StatusUnprocessableEntity, "", errors.Wrap(err, "parse date"))
		return
	}

	if timediff.IsPast(date) {
		out.Response = resp.NewResponse(http.StatusBadRequest, "", errors.New("cannot create dues in the past"))
		return
	}

	dues, err := d.DuesRepository.FindOtherByYYYYMM(ctx, 0, date)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "find dues by yyyy mm"))
		return
	}

	if dues.Id != 0 {
		out.Response = resp.NewResponse(http.StatusBadRequest, "", errors.New("dues for this month already exist"))
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

func (d *DuesDeps) QueryDues(ctx context.Context, cursor string) (out QueryDuesOut) {
	var err error
	out.Response = resp.NewResponse(http.StatusOK, "", nil)

	fromCursor, _ := strconv.ParseInt(cursor, 10, 64)
	dues, err := d.DuesRepository.Query(ctx, fromCursor, 25)
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
			Date:      d.Date.Format("2006-01-02"),
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
		out.Response = resp.NewResponse(http.StatusBadRequest, "", errors.Wrap(err, "parse uint"))
		return
	}

	if err = ValidateEditDuesIn(in); err != nil {
		out.Response = resp.NewResponse(http.StatusUnprocessableEntity, "", errors.Wrap(err, "edit dues validation"))
		return
	}

	date, err := time.Parse("2006-01-02", in.Date)
	if err != nil {
		out.Response = resp.NewResponse(http.StatusUnprocessableEntity, "", errors.Wrap(err, "parse date"))
		return
	}

	if timediff.IsPast(date) {
		out.Response = resp.NewResponse(http.StatusBadRequest, "", errors.New("cannot edit dues to the past"))
		return
	}

	dues, err := d.DuesRepository.FindById(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		out.Response = resp.NewResponse(http.StatusNotFound, "", errors.Wrap(err, "no row find dues by id"))
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
		out.Response = resp.NewResponse(http.StatusBadRequest, "", errors.New("dues for this month already exist"))
		return
	}

	duesList, err := d.MemberDuesRepository.CheckSomeonePaid(ctx, dues.Id)
	if err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "check someone paid dues"))
		return
	}

	if len(duesList) != 0 {
		out.Response = resp.NewResponse(http.StatusBadRequest, "", errors.New("someone already paid dues, cannot change"))
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
		out.Response = resp.NewResponse(http.StatusBadRequest, "", errors.Wrap(err, "parse uint"))
		return
	}

	dues, err := d.DuesRepository.FindById(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		out.Response = resp.NewResponse(http.StatusNotFound, "", errors.Wrap(err, "no row find dues by id"))
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
		out.Response = resp.NewResponse(http.StatusBadRequest, "", errors.New("someone already paid dues, cannot change"))
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
		out.StatusCode = http.StatusBadRequest
		err = errors.Wrap(err, "parse uint")
		return
	}

	dues, err := d.DuesRepository.FindById(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		out.StatusCode = http.StatusNotFound
		err = errors.Wrap(err, "no row find dues by id")
		return
	}
	if err != nil {
		out.StatusCode = http.StatusInternalServerError
		err = errors.Wrap(err, "find dues by id")
		return
	}

	duesList, err := d.MemberDuesRepository.CheckSomeonePaid(ctx, dues.Id)
	if err != nil {
		out.StatusCode = http.StatusInternalServerError
		err = errors.Wrap(err, "check someone paid dues")
		return
	}

	isPaid := false
	if len(duesList) != 0 {
		isPaid = true
	}

	out.Res.IsPaid = isPaid

	return
}
