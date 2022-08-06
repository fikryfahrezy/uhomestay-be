package homestay

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/filetype"
	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/httpdecode"
	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/resp"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
)

var (
	ErrHomestayImageNotFound = errors.New("foto atau gambar tidak ditemukan")
	ErrNotValidHomestayImage = errors.New("file bukan bukan bertipe foto atau gambar")
)

type (
	AddHomestayImageIn struct {
		File httpdecode.FileHeader `mapstructure:"file"`
	}
	AddHomestayImageRes struct {
		Id  int64  `json:"id"`
		Url string `json:"url"`
	}
	AddHomestayImageOut struct {
		resp.Response
		Res AddHomestayImageRes
	}
)

func (d *HomestayDeps) AddHomestayImage(ctx context.Context, in AddHomestayImageIn) (out AddHomestayImageOut) {
	var err error
	out.Response = resp.NewResponse(http.StatusCreated, "", nil)

	if err = ValidateAddHomestayImageIn(in); err != nil {
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
		buff := bytes.NewBuffer(nil)
		if _, err = io.Copy(buff, file); err != nil {
			out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "read file buffer"))
			return
		}

		fileCt := http.DetectContentType(buff.Bytes())
		fmt.Println(fileCt)
		if !filetype.IsTypeAllowed(fileCt) {
			out.Response = resp.NewResponse(http.StatusUnprocessableEntity, "", ErrNotValidHomestayImage)
			return
		}

		filename := strconv.FormatInt(time.Now().Unix(), 10) + "-" + strings.Trim(in.File.Filename, " ")
		if fileUrl, err = d.Upload(filename, buff); err != nil {
			out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "upload file"))
			return
		}
	}

	re := regexp.MustCompile(`[^a-zA-Z0-9]`)
	image := HomestayImageModel{
		Name:        in.File.Filename,
		AlphnumName: string(re.ReplaceAll([]byte(in.File.Filename), []byte(" "))),
		Url:         fileUrl,
	}

	if image, err = d.HomestayImageRepository.Save(ctx, image); err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "save image"))
		return
	}

	out.Res = AddHomestayImageRes{
		Id:  int64(image.Id),
		Url: fileUrl,
	}

	return
}

type (
	RemoveHomestayImageRes struct {
		Id int64 `json:"id"`
	}
	RemoveHomestayImageOut struct {
		resp.Response
		Res RemoveImageRes
	}
)

func (d *HomestayDeps) RemoveHomestayImage(ctx context.Context, pid string) (out RemoveHomestayImageOut) {
	var err error
	out.Response = resp.NewResponse(http.StatusOK, "", nil)

	id, err := strconv.ParseUint(pid, 10, 64)
	if err != nil {
		out.Response = resp.NewResponse(http.StatusNotFound, "", ErrHomestayImageNotFound)
		return
	}

	_, err = d.HomestayImageRepository.FindById(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		out.Response = resp.NewResponse(http.StatusNotFound, "", ErrHomestayImageNotFound)
		return
	}
	if err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "find image by id"))
		return
	}

	if err = d.HomestayImageRepository.DeleteById(ctx, id); err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "delete image by id"))
		return
	}

	out.Res.Id = int64(id)

	return
}
