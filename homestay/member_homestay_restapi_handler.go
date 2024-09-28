package homestay

import (
	"encoding/json"
	"net/http"

	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/resp"
	"github.com/go-chi/chi/v5"
)

func (d *HomestayDeps) PostMemberHomestay(w http.ResponseWriter, r *http.Request) {
	uid := chi.URLParam(r, "uid")
	decoder := json.NewDecoder(r.Body)

	var in AddMemberHomestayIn
	err := decoder.Decode(&in)
	if err != nil {
		resp.NewResponse(http.StatusInternalServerError, "", err).HttpJSON(w, nil)
		return
	}

	out := d.AddMemberHomestay(r.Context(), uid, in)
	out.HttpJSON(w, resp.NewHttpBody(out.Res))
}

func (d *HomestayDeps) DeleteMemberHomestay(w http.ResponseWriter, r *http.Request) {
	uid := chi.URLParam(r, "uid")
	id := chi.URLParam(r, "id")
	out := d.RemoveMemberHomestay(r.Context(), id, uid)
	out.HttpJSON(w, resp.NewHttpBody(out.Res))
}

func (d *HomestayDeps) GetMemberHomestays(w http.ResponseWriter, r *http.Request) {
	uid := chi.URLParam(r, "uid")
	cursor := r.URL.Query().Get("cursor")
	limit := r.URL.Query().Get("limit")
	out := d.QueryMemberHomestays(r.Context(), uid, cursor, limit)
	out.HttpJSON(w, resp.NewHttpBody(out.Res))
}

func (d *HomestayDeps) GetMemberHomestay(w http.ResponseWriter, r *http.Request) {
	uid := chi.URLParam(r, "uid")
	id := chi.URLParam(r, "id")
	out := d.FindMemberHomestay(r.Context(), id, uid)
	out.HttpJSON(w, resp.NewHttpBody(out.Res))
}

func (d *HomestayDeps) PutMemberHomestay(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	uid := chi.URLParam(r, "uid")
	id := chi.URLParam(r, "id")
	var in EditMemberHomestayIn
	err := decoder.Decode(&in)
	if err != nil {
		resp.NewResponse(http.StatusInternalServerError, "", err).HttpJSON(w, nil)
		return
	}

	out := d.EditMemberHomestay(r.Context(), id, uid, in)
	out.HttpJSON(w, resp.NewHttpBody(out.Res))
}
