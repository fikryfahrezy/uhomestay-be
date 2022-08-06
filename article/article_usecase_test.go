package article_test

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/article"
	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/httpdecode"
)

func TestAddArticle(t *testing.T) {
	err := ClearTables(postgrePool)
	if err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		Name               string
		ExpectedStatusCode int
		init               func()
		In                 article.AddArticleIn
	}{
		{
			Name:               "Add Article without Img Success",
			ExpectedStatusCode: http.StatusCreated,
			init:               func() {},
			In: article.AddArticleIn{
				Title:        "Title",
				ShortDesc:    "Short Desc",
				Slug:         "slug",
				ThumbnailUrl: "",
				Content:      `{"test": "hi"}`,
			},
		},
		{
			Name:               "Add Article with Img Success",
			ExpectedStatusCode: http.StatusCreated,
			init: func() {
				imgId := "balbla/images.jpg"
				imgUrl := "http://localhost/balbla/images.jpg.jpg"
				articleRepository.SetImgUrlCache(context.Background(), imgId, imgUrl)
			},
			In: article.AddArticleIn{
				Title:        "Title",
				ShortDesc:    "Short Desc",
				Slug:         "slug",
				ThumbnailUrl: "",
				Content: fmt.Sprintf(
					`{"img": "%s"}`,
					"http://localhost/balbla/images.jpg.jpg",
				),
				ContentText: "",
			},
		},
		{
			Name:               "Add Article with Imgs Success",
			ExpectedStatusCode: http.StatusCreated,
			init: func() {
				imgId := "balbla/images.jpg"
				imgUrl := "http://localhost/balbla/images.jpg.jpg"
				articleRepository.SetImgUrlCache(context.Background(), imgId, imgUrl)

				imgId2 := "blabla/images2.jpg"
				imgUrl2 := "http://localhost/blabla/images2.jpg.jpg"
				articleRepository.SetImgUrlCache(context.Background(), imgId2, imgUrl2)
			},
			In: article.AddArticleIn{
				Title:        "Title",
				ShortDesc:    "Short Desc",
				Slug:         "slug",
				ThumbnailUrl: "",
				Content: fmt.Sprintf(
					`{"img1": "%s", "img2": "%s"}`,
					"http://localhost/balbla/images.jpg.jpg",
					"http://localhost/blabla/images2.jpg.jpg",
				),
				ContentText: "",
			},
		},
		{
			Name:               "Add Article with Img and Thumbnail Success",
			ExpectedStatusCode: http.StatusCreated,
			init: func() {
				thmId := "balbla/thm.jpg"
				thmUrl := "http://localhost/balbla/thm.jpg.jpg"
				articleRepository.SetImgUrlCache(context.Background(), thmId, thmUrl)

				imgId := "balbla/images.jpg"
				imgUrl := "http://localhost/balbla/images.jpg.jpg"
				articleRepository.SetImgUrlCache(context.Background(), imgId, imgUrl)
			},
			In: article.AddArticleIn{
				Title:        "Title",
				ShortDesc:    "Short Desc",
				Slug:         "slug",
				ThumbnailUrl: "http://localhost/balbla/thm.jpg.jpg",
				Content: fmt.Sprintf(
					`{"img": "%s"}`,
					"http://localhost/balbla/images.jpg.jpg",
				),
				ContentText: "",
			},
		},
		{
			Name:               "Add Article with Imgs and Thumbnail Success",
			ExpectedStatusCode: http.StatusCreated,
			init: func() {
				thmId := "balbla/thm.jpg"
				thmUrl := "http://localhost/balbla/thm.jpg.jpg"
				articleRepository.SetImgUrlCache(context.Background(), thmId, thmUrl)

				imgId := "balbla/images.jpg"
				imgUrl := "http://localhost/balbla/images.jpg.jpg"
				articleRepository.SetImgUrlCache(context.Background(), imgId, imgUrl)

				imgId2 := "blabla/images2.jpg"
				imgUrl2 := "http://localhost/blabla/images2.jpg.jpg"
				articleRepository.SetImgUrlCache(context.Background(), imgId2, imgUrl2)
			},
			In: article.AddArticleIn{
				Title:        "Title",
				ShortDesc:    "Short Desc",
				Slug:         "slug",
				ThumbnailUrl: "http://localhost/balbla/thm.jpg.jpg",
				Content: fmt.Sprintf(
					`{"img1": "%s", "img2": "%s"}`,
					"http://localhost/balbla/images.jpg.jpg",
					"http://localhost/blabla/images2.jpg.jpg",
				),
				ContentText: "",
			},
		},
		{
			Name:               "Add Article with Title 200 chars Success",
			ExpectedStatusCode: http.StatusCreated,
			init: func() {
			},
			In: article.AddArticleIn{
				Title:        strings.Repeat("a", 200),
				ShortDesc:    "Short Desc",
				Slug:         "slug",
				ThumbnailUrl: "",
				Content:      `{"test": "test"}`,
				ContentText:  "test",
			},
		},
		{
			Name:               "Add Article with Title over 200 chars Failed",
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			init: func() {
			},
			In: article.AddArticleIn{
				Title:        strings.Repeat("a", 201),
				ShortDesc:    "Short Desc",
				Slug:         "slug",
				ThumbnailUrl: "",
				Content:      `{"test": "test"}`,
				ContentText:  "test",
			},
		},
		{
			Name:               "Add Article with Short Desc 200 chars Success",
			ExpectedStatusCode: http.StatusCreated,
			init: func() {
			},
			In: article.AddArticleIn{
				Title:        "Title",
				ShortDesc:    strings.Repeat("a", 200),
				Slug:         "slug",
				ThumbnailUrl: "",
				Content:      `{"test": "test"}`,
				ContentText:  "test",
			},
		},
		{
			Name:               "Add Article with Short Desc over 200 chars Failed",
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			init: func() {
			},
			In: article.AddArticleIn{
				Title:        "Title",
				ShortDesc:    strings.Repeat("a", 201),
				Slug:         "slug",
				ThumbnailUrl: "",
				Content:      `{"test": "test"}`,
				ContentText:  "test",
			},
		},
		{
			Name:               "Add Article with Slug 200 chars Success",
			ExpectedStatusCode: http.StatusCreated,
			init: func() {
			},
			In: article.AddArticleIn{
				Title:        "Title",
				ShortDesc:    "Shor Desc",
				Slug:         strings.Repeat("a", 200),
				ThumbnailUrl: "",
				Content:      `{"test": "test"}`,
				ContentText:  "test",
			},
		},
		{
			Name:               "Add Article with Slug over 200 chars Failed",
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			init: func() {
			},
			In: article.AddArticleIn{
				Title:        "Title",
				ShortDesc:    "Shor Desc",
				Slug:         strings.Repeat("a", 201),
				ThumbnailUrl: "",
				Content:      `{"test": "test"}`,
				ContentText:  "test",
			},
		},
	}

	for _, c := range testCases {
		t.Run(c.Name, func(t *testing.T) {
			c.init()
			res := articleDeps.AddArticle(context.Background(), c.In)

			if res.StatusCode != c.ExpectedStatusCode {
				t.Logf("%#v", res)
				t.Log(err)
				t.Fatalf("Expected response code %d. Got %d\n", c.ExpectedStatusCode, res.StatusCode)
			}
		})
	}
}

func TestQueryArticle(t *testing.T) {
	err := ClearTables(postgrePool)
	if err != nil {
		t.Fatal(err)
	}

	_, err = articleRepository.Save(context.Background(), articleSeed)
	if err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		Name               string
		ExpectedStatusCode int
	}{
		{
			Name:               "Query Article Success",
			ExpectedStatusCode: http.StatusOK,
		},
	}

	for _, c := range testCases {
		t.Run(c.Name, func(t *testing.T) {
			res := articleDeps.QueryArticle(context.Background(), "", "")

			if res.StatusCode != c.ExpectedStatusCode {
				t.Logf("%#v", res)
				t.Log(err)
				t.Fatalf("Expected response code %d. Got %d\n", c.ExpectedStatusCode, res.StatusCode)
			}
		})
	}
}

func TestFindArticleById(t *testing.T) {
	err := ClearTables(postgrePool)
	if err != nil {
		t.Fatal(err)
	}

	article, err := articleRepository.Save(context.Background(), articleSeed)
	if err != nil {
		t.Fatal(err)
	}

	pid := strconv.FormatUint(article.Id, 10)

	testCases := []struct {
		Name               string
		ExpectedStatusCode int
		Id                 string
	}{
		{
			Name:               "Find Article By Id Success",
			ExpectedStatusCode: http.StatusOK,
			Id:                 pid,
		},
		{
			Name:               "Find Article By Id Fail, Article Not Found",
			ExpectedStatusCode: http.StatusNotFound,
			Id:                 "999",
		},
	}

	for _, c := range testCases {
		t.Run(c.Name, func(t *testing.T) {
			res := articleDeps.FindArticleById(context.Background(), c.Id)

			if res.StatusCode != c.ExpectedStatusCode {
				t.Logf("%#v", res)
				t.Log(err)
				t.Fatalf("Expected response code %d. Got %d\n", c.ExpectedStatusCode, res.StatusCode)
			}
		})
	}
}

func TestEditArticle(t *testing.T) {
	err := ClearTables(postgrePool)
	if err != nil {
		t.Fatal(err)
	}

	b, err := articleRepository.Save(context.Background(), articleSeed)
	if err != nil {
		t.Fatal(err)
	}

	pid := strconv.FormatUint(b.Id, 10)

	testCases := []struct {
		Name               string
		ExpectedStatusCode int
		Id                 string
		init               func()
		In                 article.EditArticleIn
	}{
		{
			Name:               "Edit Article Success",
			ExpectedStatusCode: http.StatusOK,
			Id:                 pid,
			init:               func() {},
			In: article.EditArticleIn{
				Title:        "Title",
				ShortDesc:    "Short Desc",
				ThumbnailUrl: "",
				Content:      `{"test": "hi"}`,
				ContentText:  "hi",
			},
		},
		{
			Name:               "Edit Article with Img Success",
			ExpectedStatusCode: http.StatusOK,
			Id:                 pid,
			init: func() {
				imgId := "balbla/images.jpg"
				imgUrl := "http://localhost/balbla/images.jpg.jpg"
				articleRepository.SetImgUrlCache(context.Background(), imgId, imgUrl)
			},
			In: article.EditArticleIn{
				Title:        "Title",
				ShortDesc:    "Short Desc",
				ThumbnailUrl: "",
				Content: fmt.Sprintf(
					`{"img": "%s"}`,
					"http://localhost/balbla/images.jpg.jpg",
				),
				ContentText: "",
			},
		},
		{
			Name:               "Edit Article with Imgs Success",
			ExpectedStatusCode: http.StatusOK,
			Id:                 pid,
			init: func() {
				imgId := "balbla/images.jpg"
				imgUrl := "http://localhost/balbla/images.jpg.jpg"
				articleRepository.SetImgUrlCache(context.Background(), imgId, imgUrl)

				imgId2 := "blabla/images2.jpg"
				imgUrl2 := "http://localhost/blabla/images2.jpg.jpg"
				articleRepository.SetImgUrlCache(context.Background(), imgId2, imgUrl2)
			},
			In: article.EditArticleIn{
				Title:        "Title",
				ShortDesc:    "Short Desc",
				ThumbnailUrl: "",
				Content: fmt.Sprintf(
					`{"img1": "%s", "img2": "%s"}`,
					"http://localhost/balbla/images.jpg.jpg",
					"http://localhost/blabla/images2.jpg.jpg",
				),
				ContentText: "",
			},
		},
		{
			Name:               "Edit Article with Img and Thumbnail Success",
			ExpectedStatusCode: http.StatusOK,
			Id:                 pid,
			init: func() {
				thmId := "balbla/thm.jpg"
				thmUrl := "http://localhost/balbla/thm.jpg.jpg"
				articleRepository.SetImgUrlCache(context.Background(), thmId, thmUrl)

				imgId := "balbla/images.jpg"
				imgUrl := "http://localhost/balbla/images.jpg.jpg"
				articleRepository.SetImgUrlCache(context.Background(), imgId, imgUrl)
			},
			In: article.EditArticleIn{
				Title:        "Title",
				ShortDesc:    "Short Desc",
				ThumbnailUrl: "http://localhost/balbla/thm.jpg.jpg",
				Content: fmt.Sprintf(
					`{"img": "%s"}`,
					"http://localhost/balbla/images.jpg.jpg",
				),
				ContentText: "",
			},
		},
		{
			Name:               "Edit Article with Imgs and Thumbnail Success",
			ExpectedStatusCode: http.StatusOK,
			Id:                 pid,
			init: func() {
				thmId := "balbla/thm.jpg"
				thmUrl := "http://localhost/balbla/thm.jpg.jpg"
				articleRepository.SetImgUrlCache(context.Background(), thmId, thmUrl)

				imgId := "balbla/images.jpg"
				imgUrl := "http://localhost/balbla/images.jpg.jpg"
				articleRepository.SetImgUrlCache(context.Background(), imgId, imgUrl)

				imgId2 := "blabla/images2.jpg"
				imgUrl2 := "http://localhost/blabla/images2.jpg.jpg"
				articleRepository.SetImgUrlCache(context.Background(), imgId2, imgUrl2)
			},
			In: article.EditArticleIn{
				Title:        "Title",
				ShortDesc:    "Short Desc",
				ThumbnailUrl: "http://localhost/balbla/thm.jpg.jpg",
				Content: fmt.Sprintf(
					`{"img1": "%s", "img2": "%s"}`,
					"http://localhost/balbla/images.jpg.jpg",
					"http://localhost/blabla/images2.jpg.jpg",
				),
				ContentText: "",
			},
		},
		{
			Name:               "Edit Article Fail, Article Not Found",
			ExpectedStatusCode: http.StatusNotFound,
			Id:                 "999",
			init:               func() {},
			In: article.EditArticleIn{
				Title:        "Title",
				ShortDesc:    "Short Desc",
				ThumbnailUrl: "",
				Content:      `{"test": "hi"}`,
				ContentText:  "hi",
			},
		},
		{
			Name:               "Add Article with Title 200 chars Success",
			Id:                 pid,
			ExpectedStatusCode: http.StatusOK,
			init:               func() {},
			In: article.EditArticleIn{
				Title:        strings.Repeat("a", 200),
				ShortDesc:    "Short Desc",
				ThumbnailUrl: "",
				Content:      `{"test": "test"}`,
				ContentText:  "test",
			},
		},
		{
			Name:               "Add Article with Title over 200 chars Failed",
			Id:                 pid,
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			init: func() {
			},
			In: article.EditArticleIn{
				Title:        strings.Repeat("a", 201),
				ShortDesc:    "Short Desc",
				ThumbnailUrl: "",
				Content:      `{"test": "test"}`,
				ContentText:  "test",
			},
		},
		{
			Name:               "Add Article with Short Desc 200 chars Success",
			Id:                 pid,
			ExpectedStatusCode: http.StatusOK,
			init: func() {
			},
			In: article.EditArticleIn{
				Title:        "Title",
				ShortDesc:    strings.Repeat("a", 200),
				ThumbnailUrl: "",
				Content:      `{"test": "test"}`,
				ContentText:  "test",
			},
		},
		{
			Name:               "Add Article with Short Desc over 200 chars Failed",
			Id:                 pid,
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			init: func() {
			},
			In: article.EditArticleIn{
				Title:        "Title",
				ShortDesc:    strings.Repeat("a", 201),
				ThumbnailUrl: "",
				Content:      `{"test": "test"}`,
				ContentText:  "test",
			},
		},
	}

	for _, c := range testCases {
		t.Run(c.Name, func(t *testing.T) {
			res := articleDeps.EditArticle(context.Background(), c.Id, c.In)

			if res.StatusCode != c.ExpectedStatusCode {
				t.Logf("%#v", res)
				t.Log(err)
				t.Fatalf("Expected response code %d. Got %d\n", c.ExpectedStatusCode, res.StatusCode)
			}
		})
	}
}

func TestRemoveArticle(t *testing.T) {
	err := ClearTables(postgrePool)
	if err != nil {
		t.Fatal(err)
	}

	b, err := articleRepository.Save(context.Background(), articleSeed)
	if err != nil {
		t.Fatal(err)
	}

	pid := strconv.FormatUint(b.Id, 10)

	testCases := []struct {
		Name               string
		ExpectedStatusCode int
		Id                 string
	}{
		{
			Name:               "Remove Article By Id Success",
			ExpectedStatusCode: http.StatusOK,
			Id:                 pid,
		},
		{
			Name:               "Remove Article By Id Fail, Article Not Found",
			ExpectedStatusCode: http.StatusNotFound,
			Id:                 "999",
		},
	}

	for _, c := range testCases {
		t.Run(c.Name, func(t *testing.T) {
			res := articleDeps.RemoveArticle(context.Background(), c.Id)

			if res.StatusCode != c.ExpectedStatusCode {
				t.Logf("%#v", res)
				t.Log(err)
				t.Fatalf("Expected response code %d. Got %d\n", c.ExpectedStatusCode, res.StatusCode)
			}
		})
	}
}

func TestUploadImg(t *testing.T) {
	err := ClearRedis(redisClient)
	if err != nil {
		t.Fatal(err)
	}

	f, err := os.OpenFile(fileDir, os.O_RDONLY, 0o444)
	if err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		Name               string
		ExpectedStatusCode int
		In                 article.UploadImgIn
	}{
		{
			Name:               "Upload Image Success",
			ExpectedStatusCode: http.StatusCreated,
			In: article.UploadImgIn{
				File: httpdecode.FileHeader{
					Filename: fileName,
					File:     f,
				},
			},
		},
	}

	for _, c := range testCases {
		t.Run(c.Name, func(t *testing.T) {
			res := articleDeps.UploadImg(context.Background(), c.In)

			if res.StatusCode != c.ExpectedStatusCode {
				t.Logf("%#v", res)
				t.Log(err)
				t.Fatalf("Expected response code %d. Got %d\n", c.ExpectedStatusCode, res.StatusCode)
			}
		})
	}
}
