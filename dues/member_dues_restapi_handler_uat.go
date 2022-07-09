package dues

import (
	"encoding/json"
	"net/http"

	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/httpdecode"
	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/resp"
	"github.com/go-chi/chi/v5"
)

func (d *DuesDeps) GetMemberDuesUat(w http.ResponseWriter, r *http.Request) {
	cursor := r.URL.Query().Get("cursor")
	id := chi.URLParam(r, "id")
	limit := r.URL.Query().Get("limit")
	out := d.QueryMemberDues(r.Context(), id, cursor, limit)
	if out.Error != nil {
		d.CaptureExeption(out.Error)
	}
	out.HttpJSON(w, resp.NewHttpBody(out.Res))
}

func (d *DuesDeps) GetMembersDuesUat(w http.ResponseWriter, r *http.Request) {
	cursor := r.URL.Query().Get("cursor")
	id := chi.URLParam(r, "id")
	limit := r.URL.Query().Get("limit")
	out := d.QueryMembersDues(r.Context(), id, cursor, limit)
	if out.Error != nil {
		d.CaptureExeption(out.Error)
	}
	out.HttpJSON(w, resp.NewHttpBody(out.Res))
}

func (d *DuesDeps) PostMemberDuesUat(w http.ResponseWriter, r *http.Request) {
	var in PayMemberDuesIn
	if err := httpdecode.Multipart(r, &in, 10*1024, httpdecode.MultipartToFileHookFunc); err != nil {
		d.CaptureExeption(err)
		resp.NewResponse(http.StatusInternalServerError, "", err).HttpJSON(w, nil)
		return
	}

	username := r.URL.Query().Get("u")
	m, err := d.MemberRepository.FindByUsername(username)
	if err != nil {
		d.CaptureExeption(err)
		resp.NewResponse(http.StatusInternalServerError, "", err).HttpJSON(w, nil)
		return
	}

	id := chi.URLParam(r, "id")
	out := d.PayMemberDues(r.Context(), m.Id.UUID.String(), id, in)
	if out.Error != nil {
		d.CaptureExeption(out.Error)
	}
	out.HttpJSON(w, resp.NewHttpBody(out.Res))
}

func (d *DuesDeps) PutMemberDuesUat(w http.ResponseWriter, r *http.Request) {
	var in EditMemberDuesIn
	if err := httpdecode.Multipart(r, &in, 10*1024, httpdecode.MultipartToFileHookFunc, httpdecode.BoolToNullBoolHookFunc); err != nil {
		d.CaptureExeption(err)
		resp.NewResponse(http.StatusInternalServerError, "", err).HttpJSON(w, nil)
		return
	}

	id := chi.URLParam(r, "id")
	out := d.EditMemberDues(r.Context(), id, in)
	if out.Error != nil {
		d.CaptureExeption(out.Error)
	}
	out.HttpJSON(w, resp.NewHttpBody(out.Res))
}

func (d *DuesDeps) PatchMemberDuesUat(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var in PaidMemberDuesIn
	err := decoder.Decode(&in)
	if err != nil {
		d.CaptureExeption(err)
		resp.NewResponse(http.StatusInternalServerError, "", err).HttpJSON(w, nil)
		return
	}

	id := chi.URLParam(r, "id")
	out := d.PaidMemberDues(r.Context(), id, in)
	if out.Error != nil {
		d.CaptureExeption(out.Error)
	}
	out.HttpJSON(w, resp.NewHttpBody(out.Res))
}
