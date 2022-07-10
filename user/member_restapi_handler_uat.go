package user

import (
	"encoding/json"
	"io"
	"net/http"
	"text/template"

	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/httpdecode"
	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/resp"
	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/go-chi/chi/v5"
)

func (d *UserDeps) PostRegisterMemberUat(w http.ResponseWriter, r *http.Request) {
	var in RegisterIn
	err := json.NewDecoder(r.Body).Decode(&in)
	if err != nil {
		d.CaptureExeption(err)
		resp.NewResponse(http.StatusInternalServerError, "", err).HttpJSON(w, nil)
		return
	}

	out := d.MemberRegister(r.Context(), in)
	out.HttpJSON(w, resp.NewHttpBody(out.Res))
}

func (d *UserDeps) PostMemberUat(w http.ResponseWriter, r *http.Request) {
	var in AddMemberIn
	if err := httpdecode.Multipart(r, &in, 10*1024, httpdecode.BoolToNullBoolHookFunc, httpdecode.MultipartToFileHookFunc); err != nil {
		d.CaptureExeption(err)
		resp.NewResponse(http.StatusInternalServerError, "", err).HttpJSON(w, nil)
		return
	}

	out := d.AddMember(r.Context(), in)
	out.HttpJSON(w, resp.NewHttpBody(out.Res))
}

func (d *UserDeps) PostLoginMemberUat(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var in LoginIn
	err := decoder.Decode(&in)
	if err != nil {
		d.CaptureExeption(err)
		resp.NewResponse(http.StatusInternalServerError, "", err).HttpJSON(w, nil)
		return
	}

	out := d.MemberLogin(r.Context(), in)
	out.HttpJSON(w, resp.NewHttpBody(out.Res))
}

func (d *UserDeps) PostLoginAdminUat(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var in LoginIn
	err := decoder.Decode(&in)
	if err != nil {
		d.CaptureExeption(err)
		resp.NewResponse(http.StatusInternalServerError, "", err).HttpJSON(w, nil)
		return
	}

	out := d.AdminLogin(r.Context(), in)
	out.HttpJSON(w, resp.NewHttpBody(out.Res))
}

func (d *UserDeps) PutMemberUat(w http.ResponseWriter, r *http.Request) {
	var in EditMemberIn
	if err := httpdecode.Multipart(r, &in, 10*1024, httpdecode.BoolToNullBoolHookFunc, httpdecode.MultipartToFileHookFunc); err != nil {
		d.CaptureExeption(err)
		resp.NewResponse(http.StatusInternalServerError, "", err).HttpJSON(w, nil)
		return
	}

	id := chi.URLParam(r, "id")
	out := d.EditMember(r.Context(), id, in)
	out.HttpJSON(w, resp.NewHttpBody(out.Res))
}

func (d *UserDeps) DeleteMemberUat(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	out := d.RemoveMember(r.Context(), id)
	out.HttpJSON(w, resp.NewHttpBody(out.Res))
}

func (d *UserDeps) GetMembersUat(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	cursor := r.URL.Query().Get("cursor")
	limit := r.URL.Query().Get("limit")
	out := d.QueryMember(r.Context(), q, cursor, limit)
	out.HttpJSON(w, resp.NewHttpBody(out.Res))
}

func (d *UserDeps) GetMemberUat(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	out := d.FindMemberDetail(r.Context(), id)
	out.HttpJSON(w, resp.NewHttpBody(out.Res))
}

func (d *UserDeps) PatchMemberApprovalUat(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	out := d.ApproveMember(r.Context(), id)
	out.HttpJSON(w, resp.NewHttpBody(out.Res))
}

func (d *UserDeps) PutMemberProfileUat(w http.ResponseWriter, r *http.Request) {
	var in EditMemberIn
	if err := httpdecode.Multipart(r, &in, 10*1024, httpdecode.BoolToNullBoolHookFunc, httpdecode.MultipartToFileHookFunc); err != nil {
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

	out := d.EditMember(r.Context(), m.Id.UUID.String(), in)
	out.HttpJSON(w, resp.NewHttpBody(out.Res))
}

func (d *UserDeps) GetProfileMemberUat(w http.ResponseWriter, r *http.Request) {
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

func (d *UserDeps) RegisterFormUat(w http.ResponseWriter, r *http.Request) {
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
