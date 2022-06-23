package dues

import (
	"encoding/json"
	"net/http"

	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/resp"
	"github.com/go-chi/chi/v5"
)

func (d *DuesDeps) PostDues(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var in AddDuesIn
	err := decoder.Decode(&in)
	if err != nil {
		resp.NewResponse(http.StatusInternalServerError, "", err).HttpJSON(w, nil)
		return
	}

	out := d.AddDues(r.Context(), in)
	out.HttpJSON(w, resp.NewHttpBody(out.Res))
}

func (d *DuesDeps) PutDues(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var in EditDuesIn
	err := decoder.Decode(&in)
	if err != nil {
		resp.NewResponse(http.StatusInternalServerError, "", err).HttpJSON(w, nil)
		return
	}

	id := chi.URLParam(r, "id")
	out := d.EditDues(r.Context(), id, in)
	out.HttpJSON(w, resp.NewHttpBody(out.Res))
}

func (d *DuesDeps) DeleteDues(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	out := d.RemoveDues(r.Context(), id)
	out.HttpJSON(w, resp.NewHttpBody(out.Res))
}

func (d *DuesDeps) GetDues(w http.ResponseWriter, r *http.Request) {
	cursor := r.URL.Query().Get("cursor")
	limit := r.URL.Query().Get("limit")
	out := d.QueryDues(r.Context(), cursor, limit)
	out.HttpJSON(w, resp.NewHttpBody(out.Res))
}

func (d *DuesDeps) GetPaidDues(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	out := d.CheckPaidDues(r.Context(), id)
	out.HttpJSON(w, resp.NewHttpBody(out.Res))
}
