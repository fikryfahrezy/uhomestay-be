package article

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

var ErrArticleNotFound = errors.New("article tidak ditemukan")

type ArticleIn struct {
	Title        string
	ShortDesc    string
	ThumbnailUrl string
	Content      string
	ContentText  string
	Slug         string
}

func (d *ArticleDeps) ArticleModelBuilder(ctx context.Context, in ArticleIn) (bm ArticleModel, err error) {
	iucs, err := d.ArticleRepository.GetImgUrlsCache(ctx)
	if err != nil {
		err = errors.Wrap(err, "get img urls cache")
		return ArticleModel{}, err
	}

	urls := make(map[string]string)
	for _, iuc := range iucs {
		urls[iuc.ImageId] = iuc.ImageUrl
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
			return ArticleModel{}, err
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
			return ArticleModel{}, err
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
			return ArticleModel{}, err
		}
	}

	err = d.ArticleRepository.DelImgUrlCache(ctx)
	if err != nil {
		err = errors.Wrap(err, "del img urls cache")
		return ArticleModel{}, err
	}

	bm = ArticleModel{
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
	AddArticleIn struct {
		Title        string `json:"title"`
		ShortDesc    string `json:"short_desc"`
		ThumbnailUrl string `json:"thumbnail_url"`
		Content      string `json:"content"`
		ContentText  string `json:"content_text"`
		Slug         string `json:"slug"`
	}
	AddArticleRes struct {
		Id int64 `json:"id"`
	}
	AddArticleOut struct {
		resp.Response
		Res AddArticleRes
	}
)

func (d *ArticleDeps) AddArticle(ctx context.Context, in AddArticleIn) (out AddArticleOut) {
	var err error
	out.Response = resp.NewResponse(http.StatusCreated, "", nil)

	if err = ValidateAddArticleIn(in); err != nil {
		out.Response = resp.NewResponse(http.StatusUnprocessableEntity, "", err)
		return
	}

	article, err := d.ArticleModelBuilder(ctx, ArticleIn(in))
	if err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "article model builder"))
		return
	}

	if article, err = d.ArticleRepository.Save(ctx, article); err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "save article"))
		return
	}

	out.Res.Id = int64(article.Id)

	return
}

type (
	ArticleOut struct {
		Id           int64  `json:"id"`
		Title        string `json:"title"`
		ShortDesc    string `json:"short_desc"`
		ThumbnailUrl string `json:"thumbnail_url"`
		Slug         string `json:"slug"`
		CreatedAt    string `json:"created_at"`
	}
	QueryArticleRes struct {
		Total    int64        `json:"total"`
		Cursor   int64        `json:"cursor"`
		Articles []ArticleOut `json:"articles"`
	}
	QueryArticleOut struct {
		resp.Response
		Res QueryArticleRes
	}
)

func (d *ArticleDeps) QueryArticle(ctx context.Context, q, cursor string) (out QueryArticleOut) {
	var err error
	out.Response = resp.NewResponse(http.StatusOK, "", nil)

	articleNumber, err := d.ArticleRepository.CountArticle(ctx)
	if err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "count article"))
		return
	}

	fromCursor, _ := strconv.ParseInt(cursor, 10, 64)
	articles, err := d.ArticleRepository.Query(ctx, q, fromCursor, 25)
	if err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "query articles"))
		return
	}

	bLen := len(articles)
	var nextCursor int64
	if bLen != 0 {
		nextCursor = int64(articles[bLen-1].Id)
	}

	outArticles := make([]ArticleOut, bLen)
	for i, b := range articles {
		outArticles[i] = ArticleOut{
			Id:           int64(b.Id),
			Title:        b.Title,
			ShortDesc:    b.ShortDesc,
			Slug:         b.Slug,
			ThumbnailUrl: b.ThumbnailUrl,
			CreatedAt:    b.CreatedAt.Format("2006-01-02"),
		}
	}

	out.Res = QueryArticleRes{
		Total:    articleNumber,
		Cursor:   nextCursor,
		Articles: outArticles,
	}

	return
}

type (
	ArticleRes struct {
		Id           int64  `json:"id"`
		Title        string `json:"title"`
		ShortDesc    string `json:"short_desc"`
		ThumbnailUrl string `json:"thumbnail_url"`
		Content      string `json:"content"`
		ContentText  string `json:"content_text"`
		Slug         string `json:"slug"`
		CreatedAt    string `json:"created_at"`
	}
	FindArticleOut struct {
		resp.Response
		Res ArticleRes
	}
)

func (d *ArticleDeps) FindArticleById(ctx context.Context, pid string) (out FindArticleOut) {
	var err error
	out.Response = resp.NewResponse(http.StatusOK, "", nil)

	id, err := strconv.ParseUint(pid, 10, 64)
	if err != nil {
		out.Response = resp.NewResponse(http.StatusNotFound, "", ErrArticleNotFound)
		return
	}

	article, err := d.ArticleRepository.FindUndeletedById(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		out.Response = resp.NewResponse(http.StatusNotFound, "", ErrArticleNotFound)
		return
	}
	if err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "find article by id"))
		return
	}

	b := []byte("")
	if article.Content != nil && len(article.Content) != 0 {
		b, err = json.Marshal(article.Content)
		if err != nil {
			out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "json marshal"))
			return
		}
	}

	out.Res = ArticleRes{
		Id:           int64(article.Id),
		Title:        article.Title,
		ShortDesc:    article.ShortDesc,
		Content:      string(b),
		ContentText:  article.ContentText,
		Slug:         article.Slug,
		ThumbnailUrl: article.ThumbnailUrl,
		CreatedAt:    article.CreatedAt.Format("2006-01-02"),
	}

	return
}

type (
	EditArticleIn struct {
		Title        string `json:"title"`
		ShortDesc    string `json:"short_desc"`
		ThumbnailUrl string `json:"thumbnail_url"`
		Content      string `json:"content"`
		ContentText  string `json:"content_text"`
	}
	EditArticleRes struct {
		Id int64 `json:"id"`
	}
	EditArticleOut struct {
		resp.Response
		Res EditArticleRes
	}
)

func (d *ArticleDeps) EditArticle(ctx context.Context, pid string, in EditArticleIn) (out EditArticleOut) {
	var err error
	out.Response = resp.NewResponse(http.StatusOK, "", nil)

	if err = ValidateEditArticleIn(in); err != nil {
		out.Response = resp.NewResponse(http.StatusUnprocessableEntity, "", err)
		return
	}

	id, err := strconv.ParseUint(pid, 10, 64)
	if err != nil {
		out.Response = resp.NewResponse(http.StatusNotFound, "", ErrArticleNotFound)
		return
	}

	article, err := d.ArticleRepository.FindUndeletedById(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		out.Response = resp.NewResponse(http.StatusNotFound, "", ErrArticleNotFound)
		return
	}
	if err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "find article by id"))
		return
	}

	nb, err := d.ArticleModelBuilder(ctx, ArticleIn{
		Title:        in.Title,
		ShortDesc:    in.ShortDesc,
		ThumbnailUrl: in.ThumbnailUrl,
		Content:      in.Content,
		ContentText:  in.ContentText,
	})
	if err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "article model builder"))
		return
	}

	article.Content = nb.Content
	article.Title = nb.Title
	article.ShortDesc = nb.ShortDesc
	article.ThumbnailUrl = nb.ThumbnailUrl

	if err = d.ArticleRepository.UpdateById(ctx, id, article); err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "update position by id"))
		return
	}

	out.Res.Id = int64(id)

	return
}

type (
	RemoveArticleRes struct {
		Id int64 `json:"id"`
	}
	RemoveArticleOut struct {
		resp.Response
		Res RemoveArticleRes
	}
)

func (d *ArticleDeps) RemoveArticle(ctx context.Context, pid string) (out RemoveArticleOut) {
	var err error
	out.Response = resp.NewResponse(http.StatusOK, "", nil)

	id, err := strconv.ParseUint(pid, 10, 64)
	if err != nil {
		out.Response = resp.NewResponse(http.StatusNotFound, "", ErrArticleNotFound)
		return
	}

	_, err = d.ArticleRepository.FindUndeletedById(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		out.Response = resp.NewResponse(http.StatusNotFound, "", ErrArticleNotFound)
		return
	}
	if err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "find article by id"))
		return
	}

	if err = d.ArticleRepository.DeleteById(ctx, id); err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "delete article by id"))
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

func (d *ArticleDeps) UploadImg(ctx context.Context, in UploadImgIn) (out UploadImgOut) {
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

	if err = d.ArticleRepository.SetImgUrlCache(ctx, fileId, fileUrl); err != nil {
		out.Response = resp.NewResponse(http.StatusInternalServerError, "", errors.Wrap(err, "set img url cache"))
		return
	}

	out.Res.Url = fileUrl

	return
}
