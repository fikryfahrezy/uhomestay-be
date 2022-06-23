package blog

import (
	"context"
	"encoding/json"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/httpdecode"
	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/resp"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
)

type BlogIn struct {
	Title        string
	ShortDesc    string
	ThumbnailUrl string
	Content      string
	ContentText  string
	Slug         string
}

func (d *BlogDeps) BlogModelBuilder(ctx context.Context, in BlogIn) (bm BlogModel, err error) {
	urls, err := d.BlogRepository.GetImgUrlsCache(ctx)
	if err != nil {
		err = errors.Wrap(err, "get img urls cache")
		return BlogModel{}, err
	}

	nurls := make(map[string]string)
	thumbnailUrl := in.ThumbnailUrl
	if s := strings.Split(thumbnailUrl, d.ImgCldTmpFolder); thumbnailUrl != "" && len(s) == 2 {

		so := s[1]
		i := d.ImgCldTmpFolder + strings.TrimSuffix(so, filepath.Ext(so))
		to := d.ImgClgFolder + so

		nurl, err := d.MoveFile(i, to)
		if err != nil {
			err = errors.Wrap(err, "move file")
			return BlogModel{}, err
		}

		if nurl != "" {
			thumbnailUrl = nurl
			nurls[i] = nurl
		}
	}

	for i := range urls {
		s := strings.Split(i, d.ImgCldTmpFolder)
		if len(s) != 2 {
			urls[i] = ""
			continue
		}

		to := d.ImgClgFolder + s[1]
		nurl, err := d.MoveFile(i, to)
		if err != nil {
			err = errors.Wrap(err, "move file")
			return BlogModel{}, err
		}

		if nurl != "" {
			nurls[i] = nurl
		}
	}

	nic := in.Content
	for k, u := range nurls {
		if uk := urls[k]; uk != "" {
			nic = strings.Replace(nic, urls[k], u, -1)
		}
	}

	var nc map[string]interface{}
	if nic != "" {
		if err = json.Unmarshal([]byte(nic), &nc); err != nil {
			err = errors.Wrap(err, "unmarshal json")
			return BlogModel{}, err
		}
	}

	err = d.BlogRepository.DelImgUrlCache(ctx)
	if err != nil {
		err = errors.Wrap(err, "del img urls cache")
		return BlogModel{}, err
	}

	bm = BlogModel{
		Title:        in.Title,
		ShortDesc:    in.ShortDesc,
		ThumbnailUrl: thumbnailUrl,
		Content:      nc,
		ContentText:  in.ContentText,
		Slug:         in.Slug,
	}

	return bm, nil
}

type (
	AddBlogIn struct {
		Title        string `json:"title"`
		ShortDesc    string `json:"short_desc"`
		ThumbnailUrl string `json:"thumbnail_url"`
		Content      string `json:"content"`
		ContentText  string `json:"content_text"`
		Slug         string `json:"slug"`
	}
	AddBlogRes struct {
		Id int64 `json:"id"`
	}
	AddBlogOut struct {
		resp.Response
		Res AddBlogRes
	}
)

func (d *BlogDeps) AddBlog(ctx context.Context, in AddBlogIn) (out AddBlogOut) {
	var err error
	out.Response = resp.NewResponse(http.StatusCreated, "", nil)

	if err = ValidateAddBlogIn(in); err != nil {
		out.Response = resp.NewResponse(http.StatusUnprocessableEntity, "", errors.Wrap(err, "add blog validation"))
		return
	}

	blog, err := d.BlogModelBuilder(ctx, BlogIn{
		Title:        in.Title,
		ShortDesc:    in.ShortDesc,
		ThumbnailUrl: in.ThumbnailUrl,
		Content:      in.Content,
		ContentText:  in.ContentText,
		Slug:         in.Slug,
	})
	if err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "blog model builder"))
		return
	}

	if blog, err = d.BlogRepository.Save(ctx, blog); err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "save blog"))
		return
	}

	out.Res.Id = int64(blog.Id)

	return
}

type (
	BlogOut struct {
		Id           int64  `json:"id"`
		Title        string `json:"title"`
		ShortDesc    string `json:"short_desc"`
		ThumbnailUrl string `json:"thumbnail_url"`
		Slug         string `json:"slug"`
		CreatedAt    string `json:"created_at"`
	}
	QueryBlogRes struct {
		Cursor int64     `json:"cursor"`
		Blogs  []BlogOut `json:"blogs"`
	}
	QueryBlogOut struct {
		resp.Response
		Res QueryBlogRes
	}
)

func (d *BlogDeps) QueryBlog(ctx context.Context, q, cursor string) (out QueryBlogOut) {
	var err error
	out.Response = resp.NewResponse(http.StatusOK, "", nil)

	fromCursor, _ := strconv.ParseInt(cursor, 10, 64)
	blogs, err := d.BlogRepository.Query(ctx, q, fromCursor, 25)
	if err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "query blogs"))
		return
	}

	bLen := len(blogs)
	var nextCursor int64
	if bLen != 0 {
		nextCursor = int64(blogs[bLen-1].Id)
	}

	outBlogs := make([]BlogOut, bLen)
	for i, b := range blogs {
		outBlogs[i] = BlogOut{
			Id:           int64(b.Id),
			Title:        b.Title,
			ShortDesc:    b.ShortDesc,
			Slug:         b.Slug,
			ThumbnailUrl: b.ThumbnailUrl,
			CreatedAt:    b.CreatedAt.Format("2006-01-02"),
		}
	}

	out.Res = QueryBlogRes{
		Cursor: nextCursor,
		Blogs:  outBlogs,
	}

	return
}

type (
	BlogRes struct {
		Id           int64  `json:"id"`
		Title        string `json:"title"`
		ShortDesc    string `json:"short_desc"`
		ThumbnailUrl string `json:"thumbnail_url"`
		Content      string `json:"content"`
		ContentText  string `json:"content_text"`
		Slug         string `json:"slug"`
		CreatedAt    string `json:"created_at"`
	}
	FindBlogOut struct {
		resp.Response
		Res BlogRes
	}
)

func (d *BlogDeps) FindBlogById(ctx context.Context, pid string) (out FindBlogOut) {
	var err error
	out.Response = resp.NewResponse(http.StatusOK, "", nil)

	id, err := strconv.ParseUint(pid, 10, 64)
	if err != nil {
		out.Response = resp.NewResponse(http.StatusBadRequest, "", errors.Wrap(err, "parse uint"))
		return
	}

	blog, err := d.BlogRepository.FindUndeletedById(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		out.Response = resp.NewResponse(http.StatusNotFound, "", errors.Wrap(err, "no row find blog by id"))
		return
	}
	if err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "find blog by id"))
		return
	}

	b := []byte("")
	if blog.Content != nil && len(blog.Content) != 0 {
		b, err = json.Marshal(blog.Content)
		if err != nil {
			out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "json marshal"))
			return
		}
	}

	out.Res = BlogRes{
		Id:           int64(blog.Id),
		Title:        blog.Title,
		ShortDesc:    blog.ShortDesc,
		Content:      string(b),
		ContentText:  blog.ContentText,
		Slug:         blog.Slug,
		ThumbnailUrl: blog.ThumbnailUrl,
		CreatedAt:    blog.CreatedAt.Format("2006-01-02"),
	}

	return
}

type (
	EditBlogIn struct {
		Title        string `json:"title"`
		ShortDesc    string `json:"short_desc"`
		ThumbnailUrl string `json:"thumbnail_url"`
		Content      string `json:"content"`
		ContentText  string `json:"content_text"`
	}
	EditBlogRes struct {
		Id int64 `json:"id"`
	}
	EditBlogOut struct {
		resp.Response
		Res EditBlogRes
	}
)

func (d *BlogDeps) EditBlog(ctx context.Context, pid string, in EditBlogIn) (out EditBlogOut) {
	var err error
	out.Response = resp.NewResponse(http.StatusOK, "", nil)

	if err = ValidateEditBlogIn(in); err != nil {
		out.Response = resp.NewResponse(http.StatusUnprocessableEntity, "", errors.Wrap(err, "edit blog validation"))
		return
	}

	id, err := strconv.ParseUint(pid, 10, 64)
	if err != nil {
		out.Response = resp.NewResponse(http.StatusBadRequest, "", errors.Wrap(err, "parse uint"))
		return
	}

	blog, err := d.BlogRepository.FindUndeletedById(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		out.Response = resp.NewResponse(http.StatusNotFound, "", errors.Wrap(err, "no row find blog by id"))
		return
	}
	if err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "find blog by id"))
		return
	}

	nb, err := d.BlogModelBuilder(ctx, BlogIn{
		Title:        in.Title,
		ShortDesc:    in.ShortDesc,
		ThumbnailUrl: in.ThumbnailUrl,
		Content:      in.Content,
		ContentText:  in.ContentText,
	})
	if err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "blog model builder"))
		return
	}

	blog.Content = nb.Content
	blog.Title = nb.Title
	blog.ShortDesc = nb.ShortDesc
	blog.ThumbnailUrl = nb.ThumbnailUrl

	if err = d.BlogRepository.UpdateById(ctx, id, blog); err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "update position by id"))
		return
	}

	out.Res.Id = int64(id)

	return
}

type (
	RemoveBlogRes struct {
		Id int64 `json:"id"`
	}
	RemoveBlogOut struct {
		resp.Response
		Res RemoveBlogRes
	}
)

func (d *BlogDeps) RemoveBlog(ctx context.Context, pid string) (out RemoveBlogOut) {
	var err error
	out.Response = resp.NewResponse(http.StatusOK, "", nil)

	id, err := strconv.ParseUint(pid, 10, 64)
	if err != nil {
		out.Response = resp.NewResponse(http.StatusBadRequest, "", errors.Wrap(err, "parse uint"))
		return
	}

	_, err = d.BlogRepository.FindUndeletedById(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		out.Response = resp.NewResponse(http.StatusNotFound, "", errors.Wrap(err, "no row find position by id"))
		return
	}
	if err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "find position by id"))
		return
	}

	if err = d.BlogRepository.DeleteById(ctx, id); err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "delete position by id"))
		return
	}

	out.Res.Id = int64(id)

	return
}

type (
	UploadImgIn struct {
		File httpdecode.FileHeader `mapstructure:"file"`
	}
	UploadImgRes struct {
		Url string `json:"url"`
	}
	UploadImgOut struct {
		resp.Response
		Res UploadImgRes
	}
)

func (d *BlogDeps) UploadImg(ctx context.Context, in UploadImgIn) (out UploadImgOut) {
	var err error
	out.Response = resp.NewResponse(http.StatusCreated, "", nil)

	var file httpdecode.File
	if in.File.File != nil {
		file = in.File.File
	}

	defer func() {
		if file != nil {
			file.Close()
		}
	}()

	var fileUrl, fileId string
	if file != nil {
		filename := strconv.FormatInt(time.Now().Unix(), 10) + "-" + strings.Trim(in.File.Filename, " ")
		fileUrl, fileId, err = d.Upload(filename, file)
		if err != nil {
			out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "upload file"))
			return
		}
	}

	if err = d.BlogRepository.SetImgUrlCache(ctx, fileId, fileUrl); err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "set img url cache"))
		return
	}

	out.Res.Url = fileUrl

	return
}
