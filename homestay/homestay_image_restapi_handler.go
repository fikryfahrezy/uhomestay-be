package homestay

import (
	"net/http"

	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/httpdecode"
	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/resp"
	"github.com/go-chi/chi/v5"
)

func (d *HomestayDeps) PostHomestayImage(w http.ResponseWriter, r *http.Request) {
	var in AddHomestayImageIn
	if err := httpdecode.Multipart(r, &in, 10*1024, httpdecode.IntToNulIntHookFunc, httpdecode.MultipartToFileHookFunc, httpdecode.BoolToNullBoolHookFunc); err != nil {
		resp.NewResponse(http.StatusInternalServerError, "", err).HttpJSON(w, nil)
		return
	}

	out := d.AddHomestayImage(r.Context(), in)
	out.HttpJSON(w, resp.NewHttpBody(out.Res))
}

func (d *HomestayDeps) DeleteHomestayImage(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	out := d.RemoveHomestayImage(r.Context(), id)
	out.HttpJSON(w, resp.NewHttpBody(out.Res))
}
