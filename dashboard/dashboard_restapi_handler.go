package dashboard

import (
	"net/http"

	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/resp"
)

func (d *DashboardDeps) GetPrivateDashboard(w http.ResponseWriter, r *http.Request) {
	out := d.GetPrivate(r.Context())
	out.HttpJSON(w, resp.NewHttpBody(out.Res))
}

func (d *DashboardDeps) GetPublicDashboard(w http.ResponseWriter, r *http.Request) {
	out := d.GetPublic(r.Context())
	out.HttpJSON(w, resp.NewHttpBody(out.Res))
}
