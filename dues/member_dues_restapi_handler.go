package dues

import (
	"encoding/json"
	"net/http"

	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/httpdecode"
	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/jwt"
	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/resp"
	"github.com/go-chi/chi/v5"
)

func (d *DuesDeps) GetMemberDues(w http.ResponseWriter, r *http.Request) {
	cursor := r.URL.Query().Get("cursor")
	id := chi.URLParam(r, "id")
	limit := r.URL.Query().Get("limit")
	out := d.QueryMemberDues(r.Context(), id, cursor, limit)
	out.HttpJSON(w, resp.NewHttpBody(out.Res))
}

func (d *DuesDeps) GetMembersDues(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	out := d.QueryMembersDues(r.Context(), id, QueryMembersDuesQIn{
		Cursor:    r.URL.Query().Get("cursor"),
		Limit:     r.URL.Query().Get("limit"),
		StartDate: r.URL.Query().Get("start_date"),
		EndDate:   r.URL.Query().Get("end_date"),
	})
	out.HttpJSON(w, resp.NewHttpBody(out.Res))
}

func (d *DuesDeps) PostMemberDues(w http.ResponseWriter, r *http.Request) {
	var jwtPayload jwt.JwtPrivateClaim
	if err := jwt.DecodeCustomClaims(r, &jwtPayload); err != nil {
		resp.NewResponse(http.StatusInternalServerError, "", err).HttpJSON(w, nil)
		return
	}

	var in PayMemberDuesIn
	if err := httpdecode.Multipart(r, &in, 10*1024, httpdecode.MultipartToFileHookFunc); err != nil {
		resp.NewResponse(http.StatusInternalServerError, "", err).HttpJSON(w, nil)
		return
	}

	id := chi.URLParam(r, "id")
	out := d.PayMemberDues(r.Context(), jwtPayload.Uid, id, in)
	out.HttpJSON(w, resp.NewHttpBody(out.Res))
}

func (d *DuesDeps) PutMemberDues(w http.ResponseWriter, r *http.Request) {
	var in EditMemberDuesIn
	if err := httpdecode.Multipart(r, &in, 10*1024, httpdecode.MultipartToFileHookFunc, httpdecode.BoolToNullBoolHookFunc); err != nil {
		resp.NewResponse(http.StatusInternalServerError, "", err).HttpJSON(w, nil)
		return
	}

	id := chi.URLParam(r, "id")
	out := d.EditMemberDues(r.Context(), id, in)
	out.HttpJSON(w, resp.NewHttpBody(out.Res))
}

func (d *DuesDeps) PatchMemberDues(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var in PaidMemberDuesIn
	err := decoder.Decode(&in)
	if err != nil {
		resp.NewResponse(http.StatusInternalServerError, "", err).HttpJSON(w, nil)
		return
	}

	id := chi.URLParam(r, "id")
	out := d.PaidMemberDues(r.Context(), id, in)
	out.HttpJSON(w, resp.NewHttpBody(out.Res))
}
