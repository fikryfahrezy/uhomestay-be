package dashboard

import (
	"net/http"

	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/resp"
)

func (d *DashboardDeps) GetPrivateDashboardUat(w http.ResponseWriter, r *http.Request) {
	out := d.GetPrivate(r.Context())
	if out.Error != nil {
		d.CaptureExeption(out.Error)
	}
	out.HttpJSON(w, resp.NewHttpBody(out.Res))
}

func (d *DashboardDeps) GetPublicDashboardUat(w http.ResponseWriter, r *http.Request) {
	out := d.GetPublic(r.Context())
	if out.Error != nil {
		d.CaptureExeption(out.Error)
	}
	out.HttpJSON(w, resp.NewHttpBody(out.Res))
}
