package user

import (
	"encoding/json"
	"io"
	"net/http"
	"text/template"

	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/resp"
	"github.com/go-chi/chi/v5"
)

func (d *UserDeps) PostPositionUat(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var in AddPositionIn
	err := decoder.Decode(&in)
	if err != nil {
		d.CaptureExeption(err)
		resp.NewResponse(http.StatusInternalServerError, "", err).HttpJSON(w, nil)
		return
	}

	out := d.AddPosition(r.Context(), in)
	if out.Error != nil {
		d.CaptureExeption(out.Error)
	}
	out.HttpJSON(w, resp.NewHttpBody(out.Res))
}

func (d *UserDeps) GetPositionsUat(w http.ResponseWriter, r *http.Request) {
	cursor := r.URL.Query().Get("cursor")
	limit := r.URL.Query().Get("limit")
	out := d.QueryPosition(r.Context(), cursor, limit)
	if out.Error != nil {
		d.CaptureExeption(out.Error)
	}
	out.HttpJSON(w, resp.NewHttpBody(out.Res))
}

func (d *UserDeps) GetPositionLevelsUat(w http.ResponseWriter, r *http.Request) {
	out := d.QueryPositionLevel(r.Context())
	if out.Error != nil {
		d.CaptureExeption(out.Error)
	}
	out.HttpJSON(w, resp.NewHttpBody(out.Res))
}

func (d *UserDeps) PutPositionsUat(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var in EditPositionIn
	err := decoder.Decode(&in)
	if err != nil {
		d.CaptureExeption(err)
		resp.NewResponse(http.StatusInternalServerError, "", err).HttpJSON(w, nil)
		return
	}

	id := chi.URLParam(r, "id")
	out := d.EditPosition(r.Context(), id, in)
	if out.Error != nil {
		d.CaptureExeption(out.Error)
	}
	out.HttpJSON(w, resp.NewHttpBody(out.Res))
}

func (d *UserDeps) DeletePositionUat(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	out := d.RemovePosition(r.Context(), id)
	if out.Error != nil {
		d.CaptureExeption(out.Error)
	}
	out.HttpJSON(w, resp.NewHttpBody(out.Res))
}

func (d *UserDeps) PositionFormUat(w http.ResponseWriter, r *http.Request) {
	out := d.QueryPositionLevel(r.Context())
	if out.Error != nil {
		w.WriteHeader(out.StatusCode)
		io.WriteString(w, out.Error.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "text/html")

	t, err := template.ParseFS(d.Tmpl, "tmpl/positionform.html")
	if err != nil {
		d.CaptureExeption(err)
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, err.Error())
		return
	}

	t.Execute(w, out)
}
