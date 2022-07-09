package cashflow

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/httpdecode"
	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/resp"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
)

var (
	ErrCashflowNotFound = errors.New("cashflow tidak ditemukan")
	ErrDateFormat       = errors.New("format tanggal tidak sesuai <tahun>-<bulan>-<hari>")
)

type (
	AddCashflowIn struct {
		Date      string                `mapstructure:"date"`
		IdrAmount string                `mapstructure:"idr_amount"`
		Type      string                `mapstructure:"type"`
		Note      string                `mapstructure:"note"`
		File      httpdecode.FileHeader `mapstructure:"file"`
	}
	AddCashflowRes struct {
		Id int64 `json:"id"`
	}
	AddCashflowOut struct {
		resp.Response
		Res AddCashflowRes
	}
)

func (d *CashflowDeps) AddCashflow(ctx context.Context, in AddCashflowIn) (out AddCashflowOut) {
	var err error
	out.Response = resp.NewResponse(http.StatusCreated, "", nil)

	ct, _ := typeFromString(in.Type)

	if err = ValidateAddCashflowIn(in, ct); err != nil {
		out.Response = resp.NewResponse(http.StatusUnprocessableEntity, "", err)
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

	date, err := time.Parse("2006-01-02", in.Date)
	if err != nil {
		out.Response = resp.NewResponse(http.StatusUnprocessableEntity, "", ErrDateFormat)
		return
	}

	cashflow := CashflowModel{
		Date:         date,
		IdrAmount:    in.IdrAmount,
		Type:         ct,
		Note:         in.Note,
		ProveFileUrl: fileUrl,
	}

	if cashflow, err = d.CashflowRepository.Save(ctx, cashflow); err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "save cashflow"))
		return
	}

	out.Res.Id = int64(cashflow.Id)

	return
}

type (
	CashflowOut struct {
		Id           int64  `json:"id"`
		Date         string `json:"date"`
		Note         string `json:"note"`
		Type         string `json:"type"`
		IdrAmout     string `json:"idr_amount"`
		ProveFileUrl string `json:"prove_file_url"`
	}
	CashflowRes struct {
		TotalCash   string        `json:"total_cash"`
		IncomeCash  string        `json:"income_cash"`
		OutcomeCash string        `json:"outcome_cash"`
		Cursor      int64         `json:"curor"`
		Cashflows   []CashflowOut `json:"cashflows"`
	}
	QueryCashflowOut struct {
		resp.Response
		Res CashflowRes
	}
)

func (d *CashflowDeps) QueryCashflow(ctx context.Context, cursor, limit string) (out QueryCashflowOut) {
	out.Response = resp.NewResponse(http.StatusOK, "", nil)

	duefFlow := func(ctx context.Context, status CashflowType, flow chan float64, res chan resp.Response) {
		var r resp.Response
		amts, err := d.CashflowRepository.QueryAmtByStatus(ctx, status)
		if err != nil {
			r = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "query "+status.String+" flow"))
		}

		var namts float64
		for _, u := range amts {
			cash, _ := strconv.ParseFloat(u, 64)
			namts += cash
		}

		flow <- namts
		res <- r
	}

	inFlow := make(chan float64)
	inRes := make(chan resp.Response)
	go duefFlow(ctx, Income, inFlow, inRes)

	outFlow := make(chan float64)
	outRes := make(chan resp.Response)
	go duefFlow(ctx, Outcome, outFlow, outRes)

	outCashflows := make(chan []CashflowOut)
	nextCursor := make(chan int64)
	cRes := make(chan resp.Response)

	go func(ctx context.Context, cursor, limit string, oc chan []CashflowOut, nc chan int64, res chan resp.Response) {
		var r resp.Response
		fromCursor, _ := strconv.ParseInt(cursor, 10, 64)
		nlimit, _ := strconv.ParseInt(limit, 10, 64)
		if nlimit == 0 {
			nlimit = 25
		}

		cashflows, err := d.CashflowRepository.Query(ctx, fromCursor, nlimit)
		if err != nil {
			r = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "query cashflows"))
		}
		cLen := len(cashflows)

		var nextCursor int64
		if cLen != 0 {
			nextCursor = int64(cashflows[cLen-1].Id)
		}

		outCashflows := make([]CashflowOut, cLen)
		for i, c := range cashflows {
			outCashflows[i] = CashflowOut{
				Id:           int64(c.Id),
				Date:         c.Date.Format("2006-01-02"),
				Note:         c.Note,
				Type:         c.Type.String,
				IdrAmout:     c.IdrAmount,
				ProveFileUrl: c.ProveFileUrl,
			}
		}

		oc <- outCashflows
		nc <- nextCursor
		res <- r
	}(ctx, cursor, limit, outCashflows, nextCursor, cRes)

	inFV := <-inFlow
	inRV := <-inRes
	oFV := <-outFlow
	oRV := <-outRes
	ocV := <-outCashflows
	ncV := <-nextCursor
	cRV := <-cRes

	if inRV.Error != nil {
		out.Response = inRV
		return
	}

	if oRV.Error != nil || cRV.Error != nil {
		out.Response = oRV
		return
	}

	if cRV.Error != nil {
		out.Response = cRV
		return
	}

	out.Res = CashflowRes{
		TotalCash:   strconv.FormatFloat(inFV-oFV, 'f', -1, 64),
		IncomeCash:  strconv.FormatFloat(inFV, 'f', -1, 64),
		OutcomeCash: strconv.FormatFloat(oFV, 'f', -1, 64),
		Cursor:      ncV,
		Cashflows:   ocV,
	}

	return
}

type (
	EditCashflowIn struct {
		Date      string                `mapstructure:"date"`
		IdrAmount string                `mapstructure:"idr_amount"`
		Type      string                `mapstructure:"type"`
		Note      string                `mapstructure:"note"`
		File      httpdecode.FileHeader `mapstructure:"file"`
	}
	EditCashflowRes struct {
		Id int64 `json:"id"`
	}
	EditCashflowOut struct {
		resp.Response
		Res EditCashflowRes
	}
)

func (d *CashflowDeps) EditCashflow(ctx context.Context, pid string, in EditCashflowIn) (out EditCashflowOut) {
	var err error
	out.Response = resp.NewResponse(http.StatusOK, "", nil)

	id, err := strconv.ParseUint(pid, 10, 64)
	if err != nil {
		out.Response = resp.NewResponse(http.StatusNotFound, "", ErrCashflowNotFound)
		return
	}

	ct, _ := typeFromString(in.Type)

	if err = ValidateEditCashflowIn(in, ct); err != nil {
		out.Response = resp.NewResponse(http.StatusUnprocessableEntity, "", err)
		return
	}

	date, err := time.Parse("2006-01-02", in.Date)
	if err != nil {
		out.Response = resp.NewResponse(http.StatusUnprocessableEntity, "", ErrDateFormat)
		return
	}

	cashflow, err := d.CashflowRepository.FindById(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		out.Response = resp.NewResponse(http.StatusNotFound, "", ErrCashflowNotFound)
		return
	}
	if err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "find cashflow by id"))
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

	cashflow.Date = date
	cashflow.IdrAmount = in.IdrAmount
	cashflow.Type = ct
	cashflow.Note = in.Note

	if fileUrl != "" {
		cashflow.ProveFileUrl = fileUrl
	}

	if err = d.CashflowRepository.UpdateById(ctx, id, cashflow); err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "update cashflow by id"))
		return
	}

	out.Res.Id = int64(id)

	return
}

type (
	RemoveCashflowRes struct {
		Id int64 `json:"id"`
	}
	RemoveCashflowOut struct {
		resp.Response
		Res RemoveCashflowRes
	}
)

func (d *CashflowDeps) RemoveCashflow(ctx context.Context, pid string) (out RemoveCashflowOut) {
	var err error
	out.Response = resp.NewResponse(http.StatusOK, "", nil)

	id, err := strconv.ParseUint(pid, 10, 64)
	if err != nil {
		out.Response = resp.NewResponse(http.StatusNotFound, "", ErrCashflowNotFound)
		return
	}

	_, err = d.CashflowRepository.FindById(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		out.Response = resp.NewResponse(http.StatusNotFound, "", ErrCashflowNotFound)
		return
	}
	if err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "find cashflow by id"))
		return
	}

	if err = d.CashflowRepository.DeleteById(ctx, id); err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "delete cashflow by id"))
		return
	}

	out.Res.Id = int64(id)

	return
}
