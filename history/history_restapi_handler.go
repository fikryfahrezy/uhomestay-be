package history

import (
	"encoding/json"
	"net/http"

	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/resp"
)

func (d *HistoryDeps) PostHistory(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var in AddHistoryIn
	err := decoder.Decode(&in)
	if err != nil {
		d.CaptureExeption(err)
		resp.NewResponse(http.StatusInternalServerError, "", err).HttpJSON(w, nil)
		return
	}

	out := d.AddHistory(r.Context(), in)
	if out.Error != nil {
		d.CaptureExeption(out.Error)
	}
	out.HttpJSON(w, resp.NewHttpBody(out.Res))
}

func (d *HistoryDeps) GetHistory(w http.ResponseWriter, r *http.Request) {
	out := d.FindLatestHistory(r.Context())
	if out.Error != nil {
		d.CaptureExeption(out.Error)
	}
	out.HttpJSON(w, resp.NewHttpBody(out.Res))
}
