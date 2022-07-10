package document

import (
	"encoding/json"
	"net/http"

	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/httpdecode"
	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/resp"
	"github.com/go-chi/chi/v5"
)

func (d *DocumentDeps) PostDirDocumentUat(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var in AddDirDocumentIn
	err := decoder.Decode(&in)
	if err != nil {
		d.CaptureExeption(err)
		resp.NewResponse(http.StatusInternalServerError, "", err).HttpJSON(w, nil)
		return
	}

	out := d.AddDirDocument(r.Context(), in)
	if out.Error != nil {
		d.CaptureExeption(out.Error)
	}
	out.HttpJSON(w, resp.NewHttpBody(out.Res))
}

func (d *DocumentDeps) PostFileDocumentUat(w http.ResponseWriter, r *http.Request) {
	var in AddFileDocumentIn
	if err := httpdecode.Multipart(r, &in, 10*1024, httpdecode.IntToNulIntHookFunc, httpdecode.MultipartToFileHookFunc, httpdecode.BoolToNullBoolHookFunc); err != nil {
		d.CaptureExeption(err)
		resp.NewResponse(http.StatusInternalServerError, "", err).HttpJSON(w, nil)
		return
	}

	out := d.AddFileDocument(r.Context(), in)
	if out.Error != nil {
		d.CaptureExeption(out.Error)
	}
	out.HttpJSON(w, resp.NewHttpBody(out.Res))
}

func (d *DocumentDeps) PutDirDocumentUat(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var in EditDirDocumentIn
	err := decoder.Decode(&in)
	if err != nil {
		d.CaptureExeption(err)
		resp.NewResponse(http.StatusInternalServerError, "", err).HttpJSON(w, nil)
		return
	}

	id := chi.URLParam(r, "id")
	out := d.EditDirDocument(r.Context(), id, in)
	if out.Error != nil {
		d.CaptureExeption(out.Error)
	}
	out.HttpJSON(w, resp.NewHttpBody(out.Res))
}

func (d *DocumentDeps) PutFileDocumentUat(w http.ResponseWriter, r *http.Request) {
	var in EditFileDocumentIn
	if err := httpdecode.Multipart(r, &in, 10*1024, httpdecode.MultipartToFileHookFunc, httpdecode.BoolToNullBoolHookFunc); err != nil {
		d.CaptureExeption(err)
		resp.NewResponse(http.StatusInternalServerError, "", err).HttpJSON(w, nil)
		return
	}

	id := chi.URLParam(r, "id")
	out := d.EditFileDocument(r.Context(), id, in)
	if out.Error != nil {
		d.CaptureExeption(out.Error)
	}
	out.HttpJSON(w, resp.NewHttpBody(out.Res))
}

func (d *DocumentDeps) DeleteDocumentUat(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	out := d.RemoveDocument(r.Context(), id)
	if out.Error != nil {
		d.CaptureExeption(out.Error)
	}
	out.HttpJSON(w, resp.NewHttpBody(out.Res))
}

func (d *DocumentDeps) GetDocumentsUat(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	cursor := r.URL.Query().Get("cursor")
	limit := r.URL.Query().Get("limit")
	out := d.QueryDocument(r.Context(), q, cursor, limit)
	if out.Error != nil {
		d.CaptureExeption(out.Error)
	}
	out.HttpJSON(w, resp.NewHttpBody(out.Res))
}

func (d *DocumentDeps) GetDocumentChildrenUat(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	id := chi.URLParam(r, "id")
	cursor := r.URL.Query().Get("cursor")
	out := d.FindDocumentChildren(r.Context(), id, q, cursor)
	if out.Error != nil {
		d.CaptureExeption(out.Error)
	}
	out.HttpJSON(w, resp.NewHttpBody(out.Res))
}
