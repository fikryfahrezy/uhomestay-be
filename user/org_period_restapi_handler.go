package user

import (
	"encoding/json"
	"io"
	"net/http"
	"text/template"

	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/resp"
	"github.com/go-chi/chi/v5"
)

func (d *UserDeps) PostPeriod(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var in AddPeriodIn
	err := decoder.Decode(&in)
	if err != nil {
		resp.NewResponse(http.StatusInternalServerError, "", err).HttpJSON(w, nil)
		return
	}

	out := d.AddPeriod(r.Context(), in)
	out.HttpJSON(w, resp.NewHttpBody(out.Res))
}

func (d *UserDeps) GetPeriods(w http.ResponseWriter, r *http.Request) {
	cursor := r.URL.Query().Get("cursor")
	out := d.QueryPeriod(r.Context(), cursor)
	out.HttpJSON(w, resp.NewHttpBody(out.Res))
}

func (d *UserDeps) GetActivePeriod(w http.ResponseWriter, r *http.Request) {
	out := d.FindActivePeriod(r.Context())
	out.HttpJSON(w, resp.NewHttpBody(out.Res))
}

func (d *UserDeps) PutPeriod(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var in EditPeriodIn
	err := decoder.Decode(&in)
	if err != nil {
		resp.NewResponse(http.StatusInternalServerError, "", err).HttpJSON(w, nil)
		return
	}

	id := chi.URLParam(r, "id")
	out := d.EditPeriod(r.Context(), id, in)
	out.HttpJSON(w, resp.NewHttpBody(out.Res))
}

func (d *UserDeps) DeletePeriod(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	out := d.RemovePeriod(r.Context(), id)
	out.HttpJSON(w, resp.NewHttpBody(out.Res))
}

func (d *UserDeps) PatchPeriodStatus(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var in SwitchPeriodStatusIn
	err := decoder.Decode(&in)
	if err != nil {
		resp.NewResponse(http.StatusInternalServerError, "", err).HttpJSON(w, nil)
		return
	}

	id := chi.URLParam(r, "id")
	out := d.SwitchPeriodStatus(r.Context(), id, in)
	out.HttpJSON(w, resp.NewHttpBody(out.Res))
}

func (d *UserDeps) GetPeriodStructure(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	out := d.QueryPeriodStructure(r.Context(), id)
	out.HttpJSON(w, resp.NewHttpBody(out.Res))
}

func (d *UserDeps) PeriodForm(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "text/html")

	t, err := template.ParseFS(d.Tmpl, "tmpl/periodform.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, err.Error())
		return
	}

	t.Execute(w, nil)
}
