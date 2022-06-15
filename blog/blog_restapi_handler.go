package blog

import (
	"encoding/json"
	"net/http"

	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/httpdecode"
	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/resp"
	"github.com/go-chi/chi/v5"
)

func (d *BlogDeps) PostBlog(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var in AddBlogIn
	err := decoder.Decode(&in)
	if err != nil {
		resp.NewResponse(http.StatusInternalServerError, "", err).HttpJSON(w, nil)
		return
	}

	out := d.AddBlog(r.Context(), in)
	out.HttpJSON(w, resp.NewHttpBody(out.Res))
}

func (d *BlogDeps) GetBlogs(w http.ResponseWriter, r *http.Request) {
	cursor := r.URL.Query().Get("cursor")
	out := d.QueryBlog(r.Context(), cursor)
	out.HttpJSON(w, resp.NewHttpBody(out.Res))
}

func (d *BlogDeps) GetBlog(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	out := d.FindBlogById(r.Context(), idParam)
	out.HttpJSON(w, resp.NewHttpBody(out.Res))
}

func (d *BlogDeps) PutBlogs(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var in EditBlogIn
	err := decoder.Decode(&in)
	if err != nil {
		resp.NewResponse(http.StatusInternalServerError, "", err).HttpJSON(w, nil)
		return
	}

	idParam := chi.URLParam(r, "id")
	out := d.EditBlog(r.Context(), idParam, in)
	out.HttpJSON(w, resp.NewHttpBody(out.Res))
}

func (d *BlogDeps) DeleteBlog(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	out := d.RemoveBlog(r.Context(), idParam)
	out.HttpJSON(w, resp.NewHttpBody(out.Res))
}

func (d *BlogDeps) PostImage(w http.ResponseWriter, r *http.Request) {
	var in UploadImgIn
	if err := httpdecode.Multipart(r, &in, 10*1024, httpdecode.MultipartToFileHookFunc); err != nil {
		resp.NewResponse(http.StatusInternalServerError, "", err).HttpJSON(w, nil)
		return
	}

	out := d.UploadImg(r.Context(), in)
	out.HttpJSON(w, resp.NewHttpBody(out.Res))
}
