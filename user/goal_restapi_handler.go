package user

import (
	"encoding/json"
	"net/http"

	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/resp"
	"github.com/go-chi/chi/v5"
)

func (d *UserDeps) PostGoal(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var in AddGoalIn
	err := decoder.Decode(&in)
	if err != nil {
		resp.NewResponse(http.StatusInternalServerError, "", err).HttpJSON(w, nil)
		return
	}

	out := d.AddGoal(r.Context(), in)
	out.HttpJSON(w, resp.NewHttpBody(out.Res))
}

func (d *UserDeps) GetOrgPeriodGoal(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	out := d.FindOrgPeriodGoal(r.Context(), id)
	out.HttpJSON(w, resp.NewHttpBody(out.Res))
}
