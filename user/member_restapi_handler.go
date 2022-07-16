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
		d.CaptureExeption(err)
		resp.NewResponse(http.StatusInternalServerError, "", err).HttpJSON(w, nil)
		return
	}

	out := d.MemberRegister(r.Context(), in)
	if out.Error != nil {
		d.CaptureExeption(out.Error)
	}
	out.HttpJSON(w, resp.NewHttpBody(out.Res))
}

func (d *UserDeps) PostMember(w http.ResponseWriter, r *http.Request) {
	var in AddMemberIn
	if err := httpdecode.MultipartX(r, &in, 10*1024, httpdecode.BoolToNullBoolHookFunc, httpdecode.MultipartToFileHookFunc); err != nil {
		d.CaptureExeption(err)
		resp.NewResponse(http.StatusInternalServerError, "", err).HttpJSON(w, nil)
		return
	}

	out := d.AddMember(r.Context(), in)
	if out.Error != nil {
		d.CaptureExeption(out.Error)
	}
	out.HttpJSON(w, resp.NewHttpBody(out.Res))
}

func (d *UserDeps) PostLoginMember(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var in LoginIn
	err := decoder.Decode(&in)
	if err != nil {
		d.CaptureExeption(err)
		resp.NewResponse(http.StatusInternalServerError, "", err).HttpJSON(w, nil)
		return
	}

	out := d.MemberLogin(r.Context(), in)
	if out.Error != nil {
		d.CaptureExeption(out.Error)
	}
	out.HttpJSON(w, resp.NewHttpBody(out.Res))
}

func (d *UserDeps) PostLoginAdmin(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var in LoginIn
	err := decoder.Decode(&in)
	if err != nil {
		d.CaptureExeption(err)
		resp.NewResponse(http.StatusInternalServerError, "", err).HttpJSON(w, nil)
		return
	}

	out := d.AdminLogin(r.Context(), in)
	if out.Error != nil {
		d.CaptureExeption(out.Error)
	}
	out.HttpJSON(w, resp.NewHttpBody(out.Res))
}

func (d *UserDeps) PutMember(w http.ResponseWriter, r *http.Request) {
	var in EditMemberIn
	if err := httpdecode.MultipartX(r, &in, 10*1024, httpdecode.BoolToNullBoolHookFunc, httpdecode.MultipartToFileHookFunc); err != nil {
		d.CaptureExeption(err)
		resp.NewResponse(http.StatusInternalServerError, "", err).HttpJSON(w, nil)
		return
	}

	id := chi.URLParam(r, "id")
	out := d.EditMember(r.Context(), id, in)
	if out.Error != nil {
		d.CaptureExeption(out.Error)
	}
	out.HttpJSON(w, resp.NewHttpBody(out.Res))
}

func (d *UserDeps) DeleteMember(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	out := d.RemoveMember(r.Context(), id)
	if out.Error != nil {
		d.CaptureExeption(out.Error)
	}
	out.HttpJSON(w, resp.NewHttpBody(out.Res))
}

func (d *UserDeps) GetMembers(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	cursor := r.URL.Query().Get("cursor")
	limit := r.URL.Query().Get("limit")
	out := d.QueryMember(r.Context(), q, cursor, limit)
	if out.Error != nil {
		d.CaptureExeption(out.Error)
	}
	out.HttpJSON(w, resp.NewHttpBody(out.Res))
}

func (d *UserDeps) GetMember(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	out := d.FindMemberDetail(r.Context(), id)
	if out.Error != nil {
		d.CaptureExeption(out.Error)
	}
	out.HttpJSON(w, resp.NewHttpBody(out.Res))
}

func (d *UserDeps) PatchMemberApproval(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	out := d.ApproveMember(r.Context(), id)
	if out.Error != nil {
		d.CaptureExeption(out.Error)
	}
	out.HttpJSON(w, resp.NewHttpBody(out.Res))
}

func (d *UserDeps) PutMemberProfile(w http.ResponseWriter, r *http.Request) {
	var jwtPayload jwt.JwtPrivateClaim
	if err := jwt.DecodeCustomClaims(r, &jwtPayload); err != nil {
		d.CaptureExeption(err)
		resp.NewResponse(http.StatusInternalServerError, "", err).HttpJSON(w, nil)
		return
	}

	var in UpdateProfileIn
	if err := httpdecode.Multipart(r, &in, 10*1024, httpdecode.BoolToNullBoolHookFunc, httpdecode.MultipartToFileHookFunc); err != nil {
		d.CaptureExeption(err)
		resp.NewResponse(http.StatusInternalServerError, "", err).HttpJSON(w, nil)
		return
	}

	out := d.UpdatProfile(r.Context(), jwtPayload.Uid, in)
	if out.Error != nil {
		d.CaptureExeption(out.Error)
	}
	out.HttpJSON(w, resp.NewHttpBody(out.Res))
}

func (d *UserDeps) GetProfileMember(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	claims := r.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)

	payload, err := json.Marshal(claims)
	if err != nil {
		d.CaptureExeption(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(payload)
}

func (d *UserDeps) RegisterForm(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "text/html")

	positions := d.QueryPosition(r.Context(), "", "999")
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
		d.CaptureExeption(err)
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

func (d *UserDeps) GetMemberJwt(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	out := d.MemberLoginWithUsername(r.Context(), username)
	if out.Error != nil {
		d.CaptureExeption(out.Error)
	}
	out.HttpJSON(w, resp.NewHttpBody(out.Res))
}

func (d *UserDeps) GetAdminJwt(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	out := d.AdminLoginWithUsername(r.Context(), username)
	if out.Error != nil {
		d.CaptureExeption(out.Error)
	}
	out.HttpJSON(w, resp.NewHttpBody(out.Res))
}

func (d *UserDeps) GetUserJwt(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	out := d.UserLoginWithUsername(r.Context(), username)
	if out.Error != nil {
		d.CaptureExeption(out.Error)
	}
	out.HttpJSON(w, resp.NewHttpBody(out.Res))
}
