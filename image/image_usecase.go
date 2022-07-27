package image

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/httpdecode"
	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/resp"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
)

var ErrImageNotFound = errors.New("gambar tidak ditemukan")

type (
	AddImageIn struct {
		Description string                `mapstructure:"description"`
		File        httpdecode.FileHeader `mapstructure:"file"`
	}
	AddImageRes struct {
		Id int64 `json:"id"`
	}
	AddImageOut struct {
		resp.Response
		Res AddImageRes
	}
)

func (d *ImageDeps) AddImage(ctx context.Context, in AddImageIn) (out AddImageOut) {
	var err error
	out.Response = resp.NewResponse(http.StatusCreated, "", nil)

	if err = ValidateAddImageIn(in); err != nil {
		out.Response = resp.NewResponse(http.StatusUnprocessableEntity, "", err)
		return
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

	re := regexp.MustCompile(`[^a-zA-Z0-9]`)
	image := ImageModel{
		Name:        in.File.Filename,
		AlphnumName: string(re.ReplaceAll([]byte(in.File.Filename), []byte(" "))),
		Url:         fileUrl,
		Description: in.Description,
	}

	fmt.Println("fdskfjkdsfjdsk")
	if image, err = d.ImageRepository.Save(ctx, image); err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "save image"))
		return
	}

	out.Res.Id = int64(image.Id)

	return
}

type (
	ImageOut struct {
		Id          int64  `json:"id"`
		Name        string `json:"name"`
		Url         string `json:"url"`
		Description string `json:"description"`
	}
	QueryImageRes struct {
		Cursor int64      `json:"cursor"`
		Total  int64      `json:"total"`
		Images []ImageOut `json:"images"`
	}
	QueryImageOut struct {
		resp.Response
		Res QueryImageRes
	}
)

func (d *ImageDeps) QueryImage(ctx context.Context, cursor, limit string) (out QueryImageOut) {
	var err error
	out.Response = resp.NewResponse(http.StatusOK, "", nil)

	fromCursor, _ := strconv.ParseInt(cursor, 10, 64)
	nlimit, _ := strconv.ParseInt(limit, 10, 64)
	if nlimit == 0 {
		nlimit = 25
	}

	imageNumber, err := d.ImageRepository.CountImage(ctx)
	if err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "count image"))
		return
	}

	images, err := d.ImageRepository.Query(ctx, fromCursor, nlimit)
	if err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "query images"))
		return
	}

	imagesLen := len(images)

	var nextCursor int64
	if imagesLen != 0 {
		nextCursor = int64(images[imagesLen-1].Id)
	}

	outImages := make([]ImageOut, imagesLen)
	for i, p := range images {
		outImages[i] = ImageOut{
			Id:          int64(p.Id),
			Name:        p.Name,
			Url:         p.Url,
			Description: p.Description,
		}
	}

	out.Res = QueryImageRes{
		Cursor: nextCursor,
		Total:  imageNumber,
		Images: outImages,
	}

	return
}

type (
	RemoveImageRes struct {
		Id int64 `json:"id"`
	}
	RemoveImageOut struct {
		resp.Response
		Res RemoveImageRes
	}
)

func (d *ImageDeps) RemoveImage(ctx context.Context, pid string) (out RemoveImageOut) {
	var err error
	out.Response = resp.NewResponse(http.StatusOK, "", nil)

	id, err := strconv.ParseUint(pid, 10, 64)
	if err != nil {
		out.Response = resp.NewResponse(http.StatusNotFound, "", ErrImageNotFound)
		return
	}

	_, err = d.ImageRepository.FindById(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		out.Response = resp.NewResponse(http.StatusNotFound, "", ErrImageNotFound)
		return
	}
	if err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "find image by id"))
		return
	}

	if err = d.ImageRepository.DeleteById(ctx, id); err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "delete image by id"))
		return
	}

	out.Res.Id = int64(id)

	return
}
