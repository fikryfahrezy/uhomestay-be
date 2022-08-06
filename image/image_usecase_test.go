package image_test

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"testing"

	arbitary "github.com/PA-D3RPLA/d3if43-htt-uhomestay/arbitrary"
	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/image"
)

func TestAddImage(t *testing.T) {
	err := ClearTables(db)
	if err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		Name               string
		ExpectedStatusCode int
		In                 image.AddImageIn
	}{
		{
			Name:               "Add Image Success",
			ExpectedStatusCode: http.StatusCreated,
			In: image.AddImageIn{
				Description: "Bla Bla Bla",
				File:        generateFile(fileDir, fileName),
			},
		},
		{
			Name:               "Add Image Fail, File not image",
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			In: image.AddImageIn{
				File: generateFile("./fixture/pdf.pdf", strings.Repeat("a", 200)),
			},
		},
		{
			Name:               "Add Image Fail, File Validation Fail",
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			In:                 image.AddImageIn{},
		},
		{
			Name:               "Add Image with filename 200 chars Success",
			ExpectedStatusCode: http.StatusCreated,
			In: image.AddImageIn{
				Description: "Bla Bla Bla",
				File:        generateFile(fileDir, strings.Repeat("a", 200)),
			},
		},
		{
			Name:               "Add Image with filename over 200 chars fail",
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			In: image.AddImageIn{
				Description: "Bla Bla Bla",
				File:        generateFile(fileDir, strings.Repeat("a", 201)),
			},
		},
	}

	for _, c := range testCases {
		t.Run(c.Name, func(t *testing.T) {
			tx, err := db.Begin(context.Background())
			if err != nil {
				t.Fatal(err)
			}

			ctx := context.WithValue(context.Background(), arbitary.TrxX{}, tx)
			res := imageDeps.AddImage(ctx, c.In)
			tx.Commit(context.Background())
			tx.Rollback(context.Background())

			if res.StatusCode != c.ExpectedStatusCode {
				t.Logf("%#v", res)
				t.Log(err)
				t.Fatalf("Expected response code %d. Got %d\n", c.ExpectedStatusCode, res.StatusCode)
			}
		})
	}
}

func TestQueryDocument(t *testing.T) {
	err := ClearTables(db)
	if err != nil {
		t.Fatal(err)
	}

	_, err = imageRepository.Save(context.Background(), fileSeed)
	if err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		Name               string
		ExpectedStatusCode int
	}{
		{
			Name:               "Query Images Success",
			ExpectedStatusCode: http.StatusOK,
		},
	}

	for _, c := range testCases {
		t.Run(c.Name, func(t *testing.T) {
			tx, err := db.Begin(context.Background())
			if err != nil {
				t.Fatal(err)
			}

			ctx := context.WithValue(context.Background(), arbitary.TrxX{}, tx)
			res := imageDeps.QueryImage(ctx, "", "0")
			tx.Commit(context.Background())
			tx.Rollback(context.Background())

			if res.StatusCode != c.ExpectedStatusCode {
				t.Logf("%#v", res)
				t.Log(err)
				t.Fatalf("Expected response code %d. Got %d\n", c.ExpectedStatusCode, res.StatusCode)
			}
		})
	}
}

func TestRemoveImage(t *testing.T) {
	err := ClearTables(db)
	if err != nil {
		t.Fatal(err)
	}

	nf, err := imageRepository.Save(context.Background(), fileSeed)
	if err != nil {
		t.Fatal(err)
	}

	fid := strconv.FormatUint(nf.Id, 10)

	testCases := []struct {
		Name               string
		ExpectedStatusCode int
		Id                 string
	}{
		{
			Name:               "Remove Image Success",
			ExpectedStatusCode: http.StatusOK,
			Id:                 fid,
		},
		{
			Name:               "Remove Image Fail, Not Found",
			ExpectedStatusCode: http.StatusNotFound,
			Id:                 "999",
		},
	}

	for _, c := range testCases {
		t.Run(c.Name, func(t *testing.T) {
			tx, err := db.Begin(context.Background())
			if err != nil {
				t.Fatal(err)
			}

			ctx := context.WithValue(context.Background(), arbitary.TrxX{}, tx)
			res := imageDeps.RemoveImage(ctx, c.Id)
			tx.Commit(context.Background())
			tx.Rollback(context.Background())

			if res.StatusCode != c.ExpectedStatusCode {
				t.Logf("%#v", res)
				t.Log(err)
				t.Fatalf("Expected response code %d. Got %d\n", c.ExpectedStatusCode, res.StatusCode)
			}
		})
	}
}
