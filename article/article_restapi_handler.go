package article

import (
	"encoding/json"
	"net/http"

	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/httpdecode"
	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/resp"
	"github.com/go-chi/chi/v5"
)

func (d *ArticleDeps) PostArticle(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var in AddArticleIn
	err := decoder.Decode(&in)
	if err != nil {
		d.CaptureExeption(err)
		resp.NewResponse(http.StatusInternalServerError, "", err).HttpJSON(w, nil)
		return
	}

	out := d.AddArticle(r.Context(), in)
	if out.Error != nil {
		d.CaptureExeption(out.Error)
	}
	out.HttpJSON(w, resp.NewHttpBody(out.Res))
}

func (d *ArticleDeps) GetArticles(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	cursor := r.URL.Query().Get("cursor")
	out := d.QueryArticle(r.Context(), q, cursor)
	if out.Error != nil {
		d.CaptureExeption(out.Error)
	}
	out.HttpJSON(w, resp.NewHttpBody(out.Res))
}

func (d *ArticleDeps) GetArticle(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	out := d.FindArticleById(r.Context(), idParam)
	if out.Error != nil {
		d.CaptureExeption(out.Error)
	}
	out.HttpJSON(w, resp.NewHttpBody(out.Res))
}

func (d *ArticleDeps) PutArticle(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var in EditArticleIn
	err := decoder.Decode(&in)
	if err != nil {
		d.CaptureExeption(err)
		resp.NewResponse(http.StatusInternalServerError, "", err).HttpJSON(w, nil)
		return
	}

	idParam := chi.URLParam(r, "id")
	out := d.EditArticle(r.Context(), idParam, in)
	if out.Error != nil {
		d.CaptureExeption(out.Error)
	}
	out.HttpJSON(w, resp.NewHttpBody(out.Res))
}

func (d *ArticleDeps) DeleteArticle(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	out := d.RemoveArticle(r.Context(), idParam)
	if out.Error != nil {
		d.CaptureExeption(out.Error)
	}
	out.HttpJSON(w, resp.NewHttpBody(out.Res))
}

func (d *ArticleDeps) PostImage(w http.ResponseWriter, r *http.Request) {
	var in UploadImgIn
	if err := httpdecode.Multipart(r, &in, 10*1024, httpdecode.MultipartToFileHookFunc); err != nil {
		d.CaptureExeption(err)
		resp.NewResponse(http.StatusInternalServerError, "", err).HttpJSON(w, nil)
		return
	}

	out := d.UploadImg(r.Context(), in)
	if out.Error != nil {
		d.CaptureExeption(out.Error)
	}
	out.HttpJSON(w, resp.NewHttpBody(out.Res))
}
