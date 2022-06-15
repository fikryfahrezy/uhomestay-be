package document

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/httpdecode"
	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/resp"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"gopkg.in/guregu/null.v4"
)

type (
	AddDirDocumentIn struct {
		Name      string    `json:"name"`
		IsPrivate null.Bool `json:"is_private"`
		DirId     null.Int  `json:"dir_id"`
	}
	AddDirDocumentRes struct {
		Id int64 `json:"id"`
	}
	AddDirDocumentOut struct {
		resp.Response
		Res AddDirDocumentRes
	}
)

func (d *DocumentDeps) AddDirDocument(ctx context.Context, in AddDirDocumentIn) (out AddDirDocumentOut) {
	var err error
	out.Response = resp.NewResponse(http.StatusCreated, "", nil)

	if err = ValidateAddDirDocumentIn(in); err != nil {
		out.Response = resp.NewResponse(http.StatusUnprocessableEntity, "", errors.Wrap(err, "add document validation"))
		return
	}

	var isPrivate bool
	if in.IsPrivate.Valid {
		isPrivate = in.IsPrivate.Bool
	}

	if in.DirId.Int64 != 0 {
		doc, err := d.DocumentRepository.FindDirById(ctx, uint64(in.DirId.Int64))
		if errors.Is(err, pgx.ErrNoRows) {
			out.Response = resp.NewResponse(http.StatusNotFound, "", errors.Wrap(err, "no row find dir document by id"))
			return
		}
		if err != nil {
			out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "find document by id"))
			return
		}

		if doc.IsPrivate {
			isPrivate = doc.IsPrivate
		}
	}

	document := DocumentModel{
		Name:      in.Name,
		Type:      Dir,
		DirId:     uint64(in.DirId.Int64),
		IsPrivate: isPrivate,
	}

	if document, err = d.DocumentRepository.Save(ctx, document); err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "save document"))
		return
	}

	out.Res.Id = int64(document.Id)

	return
}

type (
	AddFileDocumentIn struct {
		IsPrivate null.Bool             `mapstructure:"is_private"`
		DirId     null.Int              `mapstructure:"dir_id"`
		File      httpdecode.FileHeader `mapstructure:"file"`
	}
	AddFileDocumentRes struct {
		Id int64 `json:"id"`
	}
	AddFileDocumentOut struct {
		resp.Response
		Res AddFileDocumentRes
	}
)

func (d *DocumentDeps) AddFileDocument(ctx context.Context, in AddFileDocumentIn) (out AddFileDocumentOut) {
	var err error
	out.Response = resp.NewResponse(http.StatusCreated, "", nil)

	if err = ValidateAddFileDocumentIn(in); err != nil {
		out.Response = resp.NewResponse(http.StatusUnprocessableEntity, "", errors.Wrap(err, "add document validation"))
		return
	}

	var isPrivate bool
	if in.IsPrivate.Valid {
		isPrivate = in.IsPrivate.Bool
	}

	if in.DirId.Int64 != 0 {
		doc, err := d.DocumentRepository.FindDirById(ctx, uint64(in.DirId.Int64))
		if errors.Is(err, pgx.ErrNoRows) {
			out.Response = resp.NewResponse(http.StatusNotFound, "", errors.Wrap(err, "no row find dir document by id"))
			return
		}
		if err != nil {
			out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "find document by id"))
			return
		}

		if doc.IsPrivate {
			isPrivate = doc.IsPrivate
		}
	}

	var file httpdecode.File
	if in.File.File != nil {
		file = in.File.File
	}

	defer func() {
		if file != nil {
			file.Close()
		}
	}()

	var fileUrl string
	if file != nil {
		filename := strconv.FormatInt(time.Now().Unix(), 10) + "-" + strings.Trim(in.File.Filename, " ")
		if fileUrl, err = d.Upload(filename, file); err != nil {
			out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "upload file"))
			return
		}
	}

	document := DocumentModel{
		Name:      in.File.Filename,
		Type:      Filetype,
		Url:       fileUrl,
		DirId:     uint64(in.DirId.Int64),
		IsPrivate: isPrivate,
	}

	if document, err = d.DocumentRepository.Save(ctx, document); err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "save document"))
		return
	}

	out.Res.Id = int64(document.Id)

	return
}

type (
	DocumentOut struct {
		IsPrivate bool   `json:"is_private"`
		Id        int64  `json:"id"`
		DirId     int64  `json:"dir_id"`
		Name      string `json:"name"`
		Type      string `json:"type"`
		Url       string `json:"url"`
	}
	QueryDocumentRes struct {
		Cursor    int64         `json:"cursor"`
		Documents []DocumentOut `json:"documents"`
	}
	QueryDocumentOut struct {
		resp.Response
		Res QueryDocumentRes
	}
)

func (d *DocumentDeps) QueryDocument(ctx context.Context, cursor string) (out QueryDocumentOut) {
	var err error
	out.Response = resp.NewResponse(http.StatusOK, "", nil)

	fromCursor, _ := strconv.ParseInt(cursor, 10, 64)
	documents, err := d.DocumentRepository.Query(ctx, fromCursor, 25)
	if err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "query documents"))
		return
	}

	docsLen := len(documents)

	var nextCursor int64
	if docsLen != 0 {
		nextCursor = int64(documents[docsLen-1].Id)
	}

	outDocuments := make([]DocumentOut, docsLen)
	for i, p := range documents {
		outDocuments[i] = DocumentOut{
			Id:        int64(p.Id),
			Name:      p.Name,
			Type:      p.Type.String,
			Url:       p.Url,
			DirId:     int64(p.DirId),
			IsPrivate: p.IsPrivate,
		}
	}

	out.Res = QueryDocumentRes{
		Cursor:    nextCursor,
		Documents: outDocuments,
	}

	return
}

type (
	EditDirDocumentIn struct {
		Name      string    `json:"name"`
		IsPrivate null.Bool `json:"is_private"`
	}
	EditDirDocumentRes struct {
		Id int64 `json:"id"`
	}
	EditDirDocumentOut struct {
		resp.Response
		Res EditDirDocumentRes
	}
)

func (d *DocumentDeps) EditDirDocument(ctx context.Context, pid string, in EditDirDocumentIn) (out EditDirDocumentOut) {
	var err error
	out.Response = resp.NewResponse(http.StatusOK, "", nil)

	id, err := strconv.ParseUint(pid, 10, 64)
	if err != nil {
		out.Response = resp.NewResponse(http.StatusBadRequest, "", errors.Wrap(err, "parse uint"))
		return
	}

	if err = ValidateEditDirDocumentIn(in); err != nil {
		out.Response = resp.NewResponse(http.StatusUnprocessableEntity, "", errors.Wrap(err, "edit document validation"))
		return
	}

	document, err := d.DocumentRepository.FindById(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) || document.Type != Dir {
		out.Response = resp.NewResponse(http.StatusNotFound, "", errors.Wrap(err, "no row find document by id"))
		return
	}
	if err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "find document by id"))
		return
	}

	var isPrivate bool
	if in.IsPrivate.Valid {
		isPrivate = in.IsPrivate.Bool
	}

	if document.DirId != 0 {
		doc, err := d.DocumentRepository.FindDirById(ctx, document.DirId)
		if errors.Is(err, pgx.ErrNoRows) {
			out.Response = resp.NewResponse(http.StatusNotFound, "", errors.Wrap(err, "no row find dir document by id"))
			return
		}
		if err != nil {
			out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "find document by id"))
			return
		}

		if doc.IsPrivate {
			isPrivate = doc.IsPrivate
		}
	}

	document.Name = in.Name
	document.IsPrivate = isPrivate

	if err = d.DocumentRepository.UpdateById(ctx, id, document); err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "update document by id"))
		return
	}

	out.Res.Id = int64(id)

	return
}

type (
	EditFileDocumentIn struct {
		IsPrivate null.Bool             `mapstructure:"is_private"`
		File      httpdecode.FileHeader `mapstructure:"file"`
	}
	EditFileDocumentRes struct {
		Id int64 `json:"id"`
	}
	EditFileDocumentOut struct {
		resp.Response
		Res EditFileDocumentRes
	}
)

func (d *DocumentDeps) EditFileDocument(ctx context.Context, pid string, in EditFileDocumentIn) (out EditFileDocumentOut) {
	var err error
	out.Response = resp.NewResponse(http.StatusOK, "", errors.Wrap(err, "update document by id"))

	id, err := strconv.ParseUint(pid, 10, 64)
	if err != nil {
		out.Response = resp.NewResponse(http.StatusBadRequest, "", errors.Wrap(err, "parse uint"))
		return
	}

	if err = ValidateEditFileDocumentIn(in); err != nil {
		out.Response = resp.NewResponse(http.StatusUnprocessableEntity, "", errors.Wrap(err, "edit document validation"))
		return
	}

	document, err := d.DocumentRepository.FindById(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) || document.Type != Filetype {
		out.Response = resp.NewResponse(http.StatusNotFound, "", errors.Wrap(err, "no row find document by id"))
		return
	}
	if err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "find document by id"))
		return
	}

	var isPrivate bool
	if in.IsPrivate.Valid {
		isPrivate = in.IsPrivate.Bool
	}

	if document.DirId != 0 {
		doc, err := d.DocumentRepository.FindDirById(ctx, document.DirId)
		if errors.Is(err, pgx.ErrNoRows) {
			out.Response = resp.NewResponse(http.StatusNotFound, "", errors.Wrap(err, "no row find dir document by id"))
			return
		}
		if err != nil {
			out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "find document by id"))
			return
		}

		if doc.IsPrivate {
			isPrivate = doc.IsPrivate
		}
	}

	var file httpdecode.File
	if in.File.File != nil {
		file = in.File.File
	}

	defer func() {
		if file != nil {
			file.Close()
		}
	}()

	var fileUrl string
	if file != nil {
		filename := strconv.FormatInt(time.Now().Unix(), 10) + "-" + strings.Trim(in.File.Filename, " ")
		if fileUrl, err = d.Upload(filename, file); err != nil {
			out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "upload file"))
			return
		}
	}

	document.Name = in.File.Filename
	document.Url = fileUrl
	document.IsPrivate = isPrivate

	if err = d.DocumentRepository.UpdateById(ctx, id, document); err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "update document by id"))
		return
	}

	out.Res.Id = int64(id)

	return
}

type (
	RemoveDocumentRes struct {
		Id int64 `json:"id"`
	}
	RemoveDocumentOut struct {
		resp.Response
		Res RemoveDocumentRes
	}
)

func (d *DocumentDeps) RemoveDocument(ctx context.Context, pid string) (out RemoveDocumentOut) {
	var err error
	out.Response = resp.NewResponse(http.StatusOK, "", nil)

	id, err := strconv.ParseUint(pid, 10, 64)
	if err != nil {
		out.Response = resp.NewResponse(http.StatusBadRequest, "", errors.Wrap(err, "parse uint"))
		return
	}

	document, err := d.DocumentRepository.FindById(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		out.Response = resp.NewResponse(http.StatusNotFound, "", errors.Wrap(err, "no row find document by id"))
		return
	}
	if err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "find document by id"))
		return
	}

	docDirIdsM := map[uint64]bool{
		document.Id: true,
	}

	// Remove child, grand child, and so on...
	docDirIds := []uint64{document.Id}
	for len(docDirIds) != 0 {
		var docDirId uint64
		docDirId, docDirIds = docDirIds[0], docDirIds[1:]
		documents, err := d.DocumentRepository.FindAllChildren(ctx, docDirId)
		if err != nil {
			out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "find document children"))
			return
		}

		for _, v := range documents {
			_, ok := docDirIdsM[v.Id]
			if !ok {
				docDirIdsM[v.Id] = true
				docDirIds = append(docDirIds, v.Id)
			}
		}
	}

	docDirIds = []uint64{}
	for k := range docDirIdsM {
		docDirIds = append(docDirIds, k)
	}

	if err = d.DocumentRepository.DeleteInId(ctx, docDirIds); err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "delete document by id"))
		return
	}

	out.Res.Id = int64(id)

	return
}

type (
	DocumentChildrenOut struct {
		resp.Response
		Res QueryDocumentRes
	}
)

func (d *DocumentDeps) FindDocumentChildren(ctx context.Context, pid, cursor string) (out DocumentChildrenOut) {
	var err error
	out.Response = resp.NewResponse(http.StatusOK, "", nil)

	id, err := strconv.ParseUint(pid, 10, 64)
	if err != nil {
		out.Response = resp.NewResponse(http.StatusBadRequest, "", errors.Wrap(err, "parse uint"))
		return
	}

	fromCursor, _ := strconv.ParseInt(cursor, 10, 64)
	documents, err := d.DocumentRepository.FindChildren(ctx, id, fromCursor, 25)
	if err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "find document children"))
		return
	}

	docsLen := len(documents)

	var nextCursor int64
	if docsLen != 0 {
		nextCursor = int64(documents[docsLen-1].Id)
	}

	outDocuments := make([]DocumentOut, docsLen)
	for i, p := range documents {
		outDocuments[i] = DocumentOut{
			Id:        int64(p.Id),
			Name:      p.Name,
			Type:      p.Type.String,
			Url:       p.Url,
			DirId:     int64(p.DirId),
			IsPrivate: p.IsPrivate,
		}
	}

	out.Res = QueryDocumentRes{
		Cursor:    nextCursor,
		Documents: outDocuments,
	}

	return
}
