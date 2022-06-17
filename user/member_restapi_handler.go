package user

import (
	"encoding/json"
	"io"
	"net/http"
	"text/template"

	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/httpdecode"
	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/jwt"
	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/resp"
	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/go-chi/chi/v5"
)

func (d *UserDeps) PostRegisterMember(w http.ResponseWriter, r *http.Request) {
	var in RegisterIn
	err := json.NewDecoder(r.Body).Decode(&in)
	if err != nil {
		resp.NewResponse(http.StatusInternalServerError, "", err).HttpJSON(w, nil)
		return
	}

	out := d.MemberRegister(r.Context(), in)
	out.HttpJSON(w, resp.NewHttpBody(out.Res))
}

func (d *UserDeps) PostMember(w http.ResponseWriter, r *http.Request) {
	var in AddMemberIn
	if err := httpdecode.Multipart(r, &in, 10*1024, httpdecode.BoolToNullBoolHookFunc, httpdecode.MultipartToFileHookFunc); err != nil {
		resp.NewResponse(http.StatusInternalServerError, "", err).HttpJSON(w, nil)
		return
	}

	out := d.AddMember(r.Context(), in)
	out.HttpJSON(w, resp.NewHttpBody(out.Res))
}

func (d *UserDeps) PostLoginMember(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var in LoginIn
	err := decoder.Decode(&in)
	if err != nil {
		resp.NewResponse(http.StatusInternalServerError, "", err).HttpJSON(w, nil)
		return
	}

	out := d.MemberLogin(r.Context(), in)
	out.HttpJSON(w, resp.NewHttpBody(out.Res))
}

func (d *UserDeps) PostLoginAdmin(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var in LoginIn
	err := decoder.Decode(&in)
	if err != nil {
		resp.NewResponse(http.StatusInternalServerError, "", err).HttpJSON(w, nil)
		return
	}

	out := d.AdminLogin(r.Context(), in)
	out.HttpJSON(w, resp.NewHttpBody(out.Res))
}

func (d *UserDeps) PutMember(w http.ResponseWriter, r *http.Request) {
	var in EditMemberIn
	if err := httpdecode.Multipart(r, &in, 10*1024, httpdecode.BoolToNullBoolHookFunc, httpdecode.MultipartToFileHookFunc); err != nil {
		resp.NewResponse(http.StatusInternalServerError, "", err).HttpJSON(w, nil)
		return
	}

	id := chi.URLParam(r, "id")
	out := d.EditMember(r.Context(), id, in)
	out.HttpJSON(w, resp.NewHttpBody(out.Res))
}

func (d *UserDeps) DeleteMember(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	out := d.RemoveMember(r.Context(), id)
	out.HttpJSON(w, resp.NewHttpBody(out.Res))
}

func (d *UserDeps) GetMembers(w http.ResponseWriter, r *http.Request) {
	cursor := r.URL.Query().Get("cursor")
	out := d.QueryMember(r.Context(), cursor)
	out.HttpJSON(w, resp.NewHttpBody(out.Res))
}

func (d *UserDeps) GetMember(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	out := d.FindMemberDetail(r.Context(), id)
	out.HttpJSON(w, resp.NewHttpBody(out.Res))
}

func (d *UserDeps) PatchMemberApproval(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	out := d.ApproveMember(r.Context(), id)
	out.HttpJSON(w, resp.NewHttpBody(out.Res))
}

func (d *UserDeps) PutMemberProfile(w http.ResponseWriter, r *http.Request) {
	var jwtPayload jwt.JwtPrivateClaim
	if err := jwt.DecodeCustomClaims(r, &jwtPayload); err != nil {
		resp.NewResponse(http.StatusInternalServerError, "", err).HttpJSON(w, nil)
		return
	}

	var in EditMemberIn
	if err := httpdecode.Multipart(r, &in, 10*1024, httpdecode.BoolToNullBoolHookFunc, httpdecode.MultipartToFileHookFunc); err != nil {
		resp.NewResponse(http.StatusInternalServerError, "", err).HttpJSON(w, nil)
		return
	}

	out := d.EditMember(r.Context(), jwtPayload.Uid, in)
	out.HttpJSON(w, resp.NewHttpBody(out.Res))
}

func (d *UserDeps) GetProfileMember(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	claims := r.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)

	payload, err := json.Marshal(claims)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(payload)
}

func (d *UserDeps) RegisterForm(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "text/html")

	positions := d.QueryPosition(r.Context(), "")
	if positions.Error != nil {
		w.WriteHeader(positions.StatusCode)
		io.WriteString(w, positions.Error.Error())
		return
	}

	periods := d.QueryPeriod(r.Context(), "")
	if periods.Error != nil {
		w.WriteHeader(periods.StatusCode)
		io.WriteString(w, periods.Error.Error())
		return
	}

	t, err := template.ParseFS(d.Tmpl, "tmpl/registerform.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, err.Error())
		return
	}

	out := struct {
		Positions, Periods interface{}
	}{
		Positions: positions.Res.Positions,
		Periods:   periods.Res.Periods,
	}

	t.Execute(w, out)
}