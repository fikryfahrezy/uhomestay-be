package homestay_test

import (
	"context"
	"net/http"
	"strconv"
	"testing"

	arbitary "github.com/PA-D3RPLA/d3if43-htt-uhomestay/arbitrary"
	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/homestay"
)

func TestAddHomestayImage(t *testing.T) {
	err := ClearTables(db)
	if err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		Name               string
		ExpectedStatusCode int
		In                 homestay.AddHomestayImageIn
	}{
		// {
		// 	Name:               "Add Image Success",
		// 	ExpectedStatusCode: http.StatusCreated,
		// 	In: homestay.AddHomestayImageIn{
		// 		File: generateFile(fileDir, fileName),
		// 	},
		// },
		{
			Name:               "Add Image Fail, File not image type",
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			In: homestay.AddHomestayImageIn{
				File: generateFile("./fixture/pdf.pdf", fileName),
			},
		},
		// {
		// 	Name:               "Add Image Fail, File Reuired",
		// 	ExpectedStatusCode: http.StatusUnprocessableEntity,
		// 	In:                 homestay.AddHomestayImageIn{},
		// },
		// {
		// 	Name:               "Add Image with filename 200 chars Success",
		// 	ExpectedStatusCode: http.StatusCreated,
		// 	In: homestay.AddHomestayImageIn{
		// 		File: generateFile(fileDir, strings.Repeat("a", 200)),
		// 	},
		// },
		// {
		// 	Name:               "Add Image with filename over 200 chars fail",
		// 	ExpectedStatusCode: http.StatusUnprocessableEntity,
		// 	In: homestay.AddHomestayImageIn{
		// 		File: generateFile(fileDir, strings.Repeat("a", 201)),
		// 	},
		// },
	}

	for _, c := range testCases {
		t.Run(c.Name, func(t *testing.T) {
			tx, err := db.Begin(context.Background())
			if err != nil {
				t.Fatal(err)
			}

			ctx := context.WithValue(context.Background(), arbitary.TrxX{}, tx)
			res := homestayDeps.AddHomestayImage(ctx, c.In)
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

func TestRemoveHomestayImage(t *testing.T) {
	err := ClearTables(db)
	if err != nil {
		t.Fatal(err)
	}

	nf, err := createHomestayImage(homestayImageRepository, fileSeed)
	if err != nil {
		t.Fatal(err)
	}

	fid := strconv.FormatInt(nf, 10)

	testCases := []struct {
		Name               string
		ExpectedStatusCode int
		Id                 string
	}{
		{
			Name:               "Remove Homestay Image Success",
			ExpectedStatusCode: http.StatusOK,
			Id:                 fid,
		},
		{
			Name:               "Remove Homestay Image Fail, Not Found",
			ExpectedStatusCode: http.StatusNotFound,
			Id:                 "",
		},
		{
			Name:               "Remove Homestay Image Fail, Not Found",
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
			res := homestayDeps.RemoveHomestayImage(ctx, c.Id)
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
