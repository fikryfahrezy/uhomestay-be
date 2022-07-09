package user

import (
	"encoding/json"
	"io"
	"net/http"
	"text/template"

	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/resp"
	"github.com/go-chi/chi/v5"
)

func (d *UserDeps) PostPeriodUat(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var in AddPeriodIn
	err := decoder.Decode(&in)
	if err != nil {
		d.CaptureExeption(err)
		resp.NewResponse(http.StatusInternalServerError, "", err).HttpJSON(w, nil)
		return
	}

	out := d.AddPeriod(r.Context(), in)
	if out.Error != nil {
		d.CaptureExeption(out.Error)
	}
	out.HttpJSON(w, resp.NewHttpBody(out.Res))
}

func (d *UserDeps) GetPeriodsUat(w http.ResponseWriter, r *http.Request) {
	cursor := r.URL.Query().Get("cursor")
	out := d.QueryPeriod(r.Context(), cursor)
	if out.Error != nil {
		d.CaptureExeption(out.Error)
	}
	out.HttpJSON(w, resp.NewHttpBody(out.Res))
}

func (d *UserDeps) GetActivePeriodUat(w http.ResponseWriter, r *http.Request) {
	out := d.FindActivePeriod(r.Context())
	if out.Error != nil {
		d.CaptureExeption(out.Error)
	}
	out.HttpJSON(w, resp.NewHttpBody(out.Res))
}

func (d *UserDeps) PutPeriodUat(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var in EditPeriodIn
	err := decoder.Decode(&in)
	if err != nil {
		d.CaptureExeption(err)
		resp.NewResponse(http.StatusInternalServerError, "", err).HttpJSON(w, nil)
		return
	}

	id := chi.URLParam(r, "id")
	out := d.EditPeriod(r.Context(), id, in)
	if out.Error != nil {
		d.CaptureExeption(out.Error)
	}
	out.HttpJSON(w, resp.NewHttpBody(out.Res))
}

func (d *UserDeps) DeletePeriodUat(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	out := d.RemovePeriod(r.Context(), id)
	if out.Error != nil {
		d.CaptureExeption(out.Error)
	}
	out.HttpJSON(w, resp.NewHttpBody(out.Res))
}

func (d *UserDeps) PatchPeriodStatusUat(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var in SwitchPeriodStatusIn
	err := decoder.Decode(&in)
	if err != nil {
		d.CaptureExeption(err)
		resp.NewResponse(http.StatusInternalServerError, "", err).HttpJSON(w, nil)
		return
	}

	id := chi.URLParam(r, "id")
	out := d.SwitchPeriodStatus(r.Context(), id, in)
	if out.Error != nil {
		d.CaptureExeption(out.Error)
	}
	out.HttpJSON(w, resp.NewHttpBody(out.Res))
}

func (d *UserDeps) GetPeriodStructureUat(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	out := d.QueryPeriodStructure(r.Context(), id)
	if out.Error != nil {
		d.CaptureExeption(out.Error)
	}
	out.HttpJSON(w, resp.NewHttpBody(out.Res))
}

func (d *UserDeps) PeriodFormUat(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "text/html")

	t, err := template.ParseFS(d.Tmpl, "tmpl/periodform.html")
	if err != nil {
		d.CaptureExeption(err)
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, err.Error())
		return
	}

	t.Execute(w, nil)
}
