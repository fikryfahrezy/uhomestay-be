package cashflow

import (
	"net/http"

	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/httpdecode"
	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/resp"
	"github.com/go-chi/chi/v5"
)

func (d *CashflowDeps) PostCashflow(w http.ResponseWriter, r *http.Request) {
	var in AddCashflowIn
	if err := httpdecode.Multipart(r, &in, 10*1024, httpdecode.IntToNulIntHookFunc, httpdecode.MultipartToFileHookFunc, httpdecode.BoolToNullBoolHookFunc); err != nil {
		resp.NewResponse(http.StatusInternalServerError, "", err).HttpJSON(w, nil)
		return
	}

	out := d.AddCashflow(r.Context(), in)
	out.HttpJSON(w, resp.NewHttpBody(out.Res))
}

func (d *CashflowDeps) PutCashflow(w http.ResponseWriter, r *http.Request) {
	var in EditCashflowIn
	if err := httpdecode.Multipart(r, &in, 10*1024, httpdecode.IntToNulIntHookFunc, httpdecode.MultipartToFileHookFunc, httpdecode.BoolToNullBoolHookFunc); err != nil {
		resp.NewResponse(http.StatusInternalServerError, "", err).HttpJSON(w, nil)
		return
	}

	id := chi.URLParam(r, "id")
	out := d.EditCashflow(r.Context(), id, in)
	out.HttpJSON(w, resp.NewHttpBody(out.Res))
}

func (d *CashflowDeps) DeleteCashflow(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	out := d.RemoveCashflow(r.Context(), id)
	out.HttpJSON(w, resp.NewHttpBody(out.Res))
}

func (d *CashflowDeps) GetCashflows(w http.ResponseWriter, r *http.Request) {
	cursor := r.URL.Query().Get("cursor")
	out := d.QueryCashflow(r.Context(), cursor)
	out.HttpJSON(w, resp.NewHttpBody(out.Res))
}
