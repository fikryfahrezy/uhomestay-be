package blog_test

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"testing"

	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/blog"
	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/httpdecode"
)

func TestAddBlog(t *testing.T) {
	err := ClearTables(postgrePool)
	if err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		Name               string
		ExpectedStatusCode int
		init               func()
		In                 blog.AddBlogIn
	}{
		{
			Name:               "Add Blog without Img Success",
			ExpectedStatusCode: http.StatusCreated,
			init:               func() {},
			In: blog.AddBlogIn{
				Title:        "Title",
				ShortDesc:    "Short Desc",
				Slug:         "slug",
				ThumbnailUrl: "",
				Content:      `{"test": "hi"}`,
			},
		},
		{
			Name:               "Add Blog with Img Success",
			ExpectedStatusCode: http.StatusCreated,
			init: func() {
				imgId := "balbla/images.jpg"
				imgUrl := "http://localhost/balbla/images.jpg.jpg"
				blogDeps.BlogRepository.SetImgUrlCache(context.Background(), imgId, imgUrl)
			},
			In: blog.AddBlogIn{
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
			Name:               "Add Blog with Imgs Success",
			ExpectedStatusCode: http.StatusCreated,
			init: func() {
				imgId := "balbla/images.jpg"
				imgUrl := "http://localhost/balbla/images.jpg.jpg"
				blogDeps.BlogRepository.SetImgUrlCache(context.Background(), imgId, imgUrl)

				imgId2 := "blabla/images2.jpg"
				imgUrl2 := "http://localhost/blabla/images2.jpg.jpg"
				blogDeps.BlogRepository.SetImgUrlCache(context.Background(), imgId2, imgUrl2)
			},
			In: blog.AddBlogIn{
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
			Name:               "Add Blog with Img and Thumbnail Success",
			ExpectedStatusCode: http.StatusCreated,
			init: func() {
				thmId := "balbla/thm.jpg"
				thmUrl := "http://localhost/balbla/thm.jpg.jpg"
				blogDeps.BlogRepository.SetImgUrlCache(context.Background(), thmId, thmUrl)

				imgId := "balbla/images.jpg"
				imgUrl := "http://localhost/balbla/images.jpg.jpg"
				blogDeps.BlogRepository.SetImgUrlCache(context.Background(), imgId, imgUrl)
			},
			In: blog.AddBlogIn{
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
			Name:               "Add Blog with Imgs and Thumbnail Success",
			ExpectedStatusCode: http.StatusCreated,
			init: func() {
				thmId := "balbla/thm.jpg"
				thmUrl := "http://localhost/balbla/thm.jpg.jpg"
				blogDeps.BlogRepository.SetImgUrlCache(context.Background(), thmId, thmUrl)

				imgId := "balbla/images.jpg"
				imgUrl := "http://localhost/balbla/images.jpg.jpg"
				blogDeps.BlogRepository.SetImgUrlCache(context.Background(), imgId, imgUrl)

				imgId2 := "blabla/images2.jpg"
				imgUrl2 := "http://localhost/blabla/images2.jpg.jpg"
				blogDeps.BlogRepository.SetImgUrlCache(context.Background(), imgId2, imgUrl2)
			},
			In: blog.AddBlogIn{
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
	}

	for _, c := range testCases {
		t.Run(c.Name, func(t *testing.T) {
			c.init()
			res := blogDeps.AddBlog(context.Background(), c.In)

			if res.StatusCode != c.ExpectedStatusCode {
				t.Logf("%#v", res)
				t.Log(err)
				t.Fatalf("Expected response code %d. Got %d\n", c.ExpectedStatusCode, res.StatusCode)
			}
		})
	}
}

func TestQueryHistory(t *testing.T) {
	err := ClearTables(postgrePool)
	if err != nil {
		t.Fatal(err)
	}

	_, err = blogRepository.Save(context.Background(), blogSeed)
	if err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		Name               string
		ExpectedStatusCode int
	}{
		{
			Name:               "Query Blog Success",
			ExpectedStatusCode: http.StatusOK,
		},
	}

	for _, c := range testCases {
		t.Run(c.Name, func(t *testing.T) {
			res := blogDeps.QueryBlog(context.Background(), "", "")

			if res.StatusCode != c.ExpectedStatusCode {
				t.Logf("%#v", res)
				t.Log(err)
				t.Fatalf("Expected response code %d. Got %d\n", c.ExpectedStatusCode, res.StatusCode)
			}
		})
	}
}

func TestFindBlogById(t *testing.T) {
	err := ClearTables(postgrePool)
	if err != nil {
		t.Fatal(err)
	}

	blog, err := blogRepository.Save(context.Background(), blogSeed)
	if err != nil {
		t.Fatal(err)
	}

	pid := strconv.FormatUint(blog.Id, 10)

	testCases := []struct {
		Name               string
		ExpectedStatusCode int
		Id                 string
	}{
		{
			Name:               "Find Blog By Id Success",
			ExpectedStatusCode: http.StatusOK,
			Id:                 pid,
		},
		{
			Name:               "Find Blog By Id Fail, Blog Not Found",
			ExpectedStatusCode: http.StatusNotFound,
			Id:                 "999",
		},
	}

	for _, c := range testCases {
		t.Run(c.Name, func(t *testing.T) {
			res := blogDeps.FindBlogById(context.Background(), c.Id)

			if res.StatusCode != c.ExpectedStatusCode {
				t.Logf("%#v", res)
				t.Log(err)
				t.Fatalf("Expected response code %d. Got %d\n", c.ExpectedStatusCode, res.StatusCode)
			}
		})
	}
}

func TestEditBlog(t *testing.T) {
	err := ClearTables(postgrePool)
	if err != nil {
		t.Fatal(err)
	}

	b, err := blogRepository.Save(context.Background(), blogSeed)
	if err != nil {
		t.Fatal(err)
	}

	pid := strconv.FormatUint(b.Id, 10)

	testCases := []struct {
		Name               string
		ExpectedStatusCode int
		Id                 string
		init               func()
		In                 blog.EditBlogIn
	}{
		{
			Name:               "Edit Blog Success",
			ExpectedStatusCode: http.StatusOK,
			Id:                 pid,
			init:               func() {},
			In: blog.EditBlogIn{
				Title:        "Title",
				ShortDesc:    "Short Desc",
				ThumbnailUrl: "",
				Content:      `{"test": "hi"}`,
				ContentText:  "hi",
			},
		},
		{
			Name:               "Edit Blog with Img Success",
			ExpectedStatusCode: http.StatusOK,
			Id:                 pid,
			init: func() {
				imgId := "balbla/images.jpg"
				imgUrl := "http://localhost/balbla/images.jpg.jpg"
				blogDeps.BlogRepository.SetImgUrlCache(context.Background(), imgId, imgUrl)
			},
			In: blog.EditBlogIn{
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
			Name:               "Edit Blog with Imgs Success",
			ExpectedStatusCode: http.StatusOK,
			Id:                 pid,
			init: func() {
				imgId := "balbla/images.jpg"
				imgUrl := "http://localhost/balbla/images.jpg.jpg"
				blogDeps.BlogRepository.SetImgUrlCache(context.Background(), imgId, imgUrl)

				imgId2 := "blabla/images2.jpg"
				imgUrl2 := "http://localhost/blabla/images2.jpg.jpg"
				blogDeps.BlogRepository.SetImgUrlCache(context.Background(), imgId2, imgUrl2)
			},
			In: blog.EditBlogIn{
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
			Name:               "Edit Blog with Img and Thumbnail Success",
			ExpectedStatusCode: http.StatusOK,
			Id:                 pid,
			init: func() {
				thmId := "balbla/thm.jpg"
				thmUrl := "http://localhost/balbla/thm.jpg.jpg"
				blogDeps.BlogRepository.SetImgUrlCache(context.Background(), thmId, thmUrl)

				imgId := "balbla/images.jpg"
				imgUrl := "http://localhost/balbla/images.jpg.jpg"
				blogDeps.BlogRepository.SetImgUrlCache(context.Background(), imgId, imgUrl)
			},
			In: blog.EditBlogIn{
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
			Name:               "Edit Blog with Imgs and Thumbnail Success",
			ExpectedStatusCode: http.StatusOK,
			Id:                 pid,
			init: func() {
				thmId := "balbla/thm.jpg"
				thmUrl := "http://localhost/balbla/thm.jpg.jpg"
				blogDeps.BlogRepository.SetImgUrlCache(context.Background(), thmId, thmUrl)

				imgId := "balbla/images.jpg"
				imgUrl := "http://localhost/balbla/images.jpg.jpg"
				blogDeps.BlogRepository.SetImgUrlCache(context.Background(), imgId, imgUrl)

				imgId2 := "blabla/images2.jpg"
				imgUrl2 := "http://localhost/blabla/images2.jpg.jpg"
				blogDeps.BlogRepository.SetImgUrlCache(context.Background(), imgId2, imgUrl2)
			},
			In: blog.EditBlogIn{
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
			Name:               "Edit Blog Fail, Blog Not Found",
			ExpectedStatusCode: http.StatusNotFound,
			Id:                 "999",
			init:               func() {},
			In: blog.EditBlogIn{
				Title:        "Title",
				ShortDesc:    "Short Desc",
				ThumbnailUrl: "",
				Content:      `{"test": "hi"}`,
				ContentText:  "hi",
			},
		},
	}

	for _, c := range testCases {
		t.Run(c.Name, func(t *testing.T) {
			res := blogDeps.EditBlog(context.Background(), c.Id, c.In)

			if res.StatusCode != c.ExpectedStatusCode {
				t.Logf("%#v", res)
				t.Log(err)
				t.Fatalf("Expected response code %d. Got %d\n", c.ExpectedStatusCode, res.StatusCode)
			}
		})
	}
}

func TestRemoveBlog(t *testing.T) {
	err := ClearTables(postgrePool)
	if err != nil {
		t.Fatal(err)
	}

	b, err := blogRepository.Save(context.Background(), blogSeed)
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
			Name:               "Remove Blog By Id Success",
			ExpectedStatusCode: http.StatusOK,
			Id:                 pid,
		},
		{
			Name:               "Remove Blog By Id Fail, Blog Not Found",
			ExpectedStatusCode: http.StatusNotFound,
			Id:                 "999",
		},
	}

	for _, c := range testCases {
		t.Run(c.Name, func(t *testing.T) {
			res := blogDeps.RemoveBlog(context.Background(), c.Id)

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
		In                 blog.UploadImgIn
	}{
		{
			Name:               "Upload Image Success",
			ExpectedStatusCode: http.StatusCreated,
			In: blog.UploadImgIn{
				File: httpdecode.FileHeader{
					Filename: fileName,
					File:     f,
				},
			},
		},
	}

	for _, c := range testCases {
		t.Run(c.Name, func(t *testing.T) {
			res := blogDeps.UploadImg(context.Background(), c.In)

			if res.StatusCode != c.ExpectedStatusCode {
				t.Logf("%#v", res)
				t.Log(err)
				t.Fatalf("Expected response code %d. Got %d\n", c.ExpectedStatusCode, res.StatusCode)
			}
		})
	}
}
