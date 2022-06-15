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
	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/pagination"
	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/resp"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
)

type BlogIn struct {
	Title        string
	ShortDesc    string
	ThumbnailUrl string
	Content      string
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
		Slug         string `json:"slug"`
	}
	AddBlogRes struct {
		Id string `json:"id"`
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
		Slug:         in.Slug,
	})
	if err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "blog model builder"))
		return
	}

	blogId, err := d.BlogRepository.Save(ctx, blog)
	if err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "save blog"))
		return
	}

	out.Res.Id = blogId

	return
}

type (
	BlogOut struct {
		Id           string `json:"id"`
		Title        string `json:"title"`
		ShortDesc    string `json:"short_desc"`
		ThumbnailUrl string `json:"thumbnail_url"`
		Content      string `json:"content"`
		Slug         string `json:"slug"`
		CreatedAt    string `json:"created_at"`
	}
	QueryBlogRes struct {
		Cursor string    `json:"cursor"`
		Blogs  []BlogOut `json:"blogs"`
	}
	QueryBlogOut struct {
		resp.Response
		Res QueryBlogRes
	}
)

func (d *BlogDeps) QueryBlog(ctx context.Context, cursor string) (out QueryBlogOut) {
	var err error
	out.Response = resp.NewResponse(http.StatusOK, "", nil)

	s, _, err := pagination.DecodeSIDCursor(cursor)
	if err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "decode sid cursor"))
		return
	}

	blogs, err := d.BlogRepository.Query(ctx, s, 25)
	if err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "query blogs"))
		return
	}

	bLen := len(blogs)
	var nextCursor string
	if bLen != 0 {
		md := blogs[bLen-1]
		nextCursor = pagination.EncodeSIDCursor(md.Id, md.CreatedAt)
	}

	outBlogs := make([]BlogOut, bLen)
	for i, b := range blogs {
		c, err := json.Marshal(b.Content)
		if err != nil {
			out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "json marshal"))
			return
		}

		outBlogs[i] = BlogOut{
			Id:           b.Id,
			Title:        b.Title,
			ShortDesc:    b.ShortDesc,
			Content:      string(c),
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
		Id           string `json:"id"`
		Title        string `json:"title"`
		ShortDesc    string `json:"short_desc"`
		ThumbnailUrl string `json:"thumbnail_url"`
		Content      string `json:"content"`
		Slug         string `json:"slug"`
		CreatedAt    string `json:"created_at"`
	}
	FindBlogOut struct {
		resp.Response
		Res BlogRes
	}
)

func (d *BlogDeps) FindBlogById(ctx context.Context, id string) (out FindBlogOut) {
	var err error
	out.Response = resp.NewResponse(http.StatusOK, "", nil)

	blog, err := d.BlogRepository.FindUndeletedById(ctx, id)
	if errors.Is(err, mongo.ErrNoDocuments) {
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
		Id:           blog.Id,
		Title:        blog.Title,
		ShortDesc:    blog.ShortDesc,
		Content:      string(b),
		Slug:         blog.Slug,
		ThumbnailUrl: blog.ThumbnailUrl,
		CreatedAt:    blog.CreatedAt.Format("2006-01-02"),
	}

	return
}

type (
	EditBlogIn struct {
		Id           string `json:"id"`
		Title        string `json:"title"`
		ShortDesc    string `json:"short_desc"`
		ThumbnailUrl string `json:"thumbnail_url"`
		Content      string `json:"content"`
		Slug         string `json:"slug"`
	}
	EditBlogRes struct {
		Id string `json:"id"`
	}
	EditBlogOut struct {
		resp.Response
		Res EditBlogRes
	}
)

func (d *BlogDeps) EditBlog(ctx context.Context, id string, in EditBlogIn) (out EditBlogOut) {
	var err error
	out.Response = resp.NewResponse(http.StatusOK, "", nil)

	if err = ValidateEditBlogIn(in); err != nil {
		out.Response = resp.NewResponse(http.StatusUnprocessableEntity, "", errors.Wrap(err, "edit blog validation"))
		return
	}

	blog, err := d.BlogRepository.FindUndeletedById(ctx, id)
	if errors.Is(err, mongo.ErrNoDocuments) {
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
		Slug:         in.Slug,
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

	out.Res.Id = id

	return
}

type (
	RemoveBlogRes struct {
		Id string `json:"id"`
	}
	RemoveBlogOut struct {
		resp.Response
		Res RemoveBlogRes
	}
)

func (d *BlogDeps) RemoveBlog(ctx context.Context, id string) (out RemoveBlogOut) {
	var err error
	out.Response = resp.NewResponse(http.StatusOK, "", nil)

	_, err = d.BlogRepository.FindUndeletedById(ctx, id)
	if errors.Is(err, mongo.ErrNoDocuments) {
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

	out.Res.Id = id

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
