package document_test

import (
	"context"
	"net/http"
	"os"
	"strconv"
	"testing"

	arbitary "github.com/PA-D3RPLA/d3if43-htt-uhomestay/arbitrary"
	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/document"
	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/httpdecode"
	"gopkg.in/guregu/null.v4"
)

func TestAddDirDocument(t *testing.T) {
	err := ClearTables(db)
	if err != nil {
		t.Fatal(err)
	}

	nd, err := documentRepository.Save(context.Background(), dirSeed)
	if err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		Name               string
		ExpectedStatusCode int
		In                 document.AddDirDocumentIn
	}{
		{
			Name:               "Add Dir Document at root dir Success",
			ExpectedStatusCode: http.StatusCreated,
			In: document.AddDirDocumentIn{
				Name:      "Dir A",
				DirId:     null.IntFrom(0),
				IsPrivate: null.BoolFrom(false),
			},
		},
		{
			Name:               "Add Dir Document at 'n' dir Success",
			ExpectedStatusCode: http.StatusCreated,
			In: document.AddDirDocumentIn{
				Name:      "Dir A",
				DirId:     null.IntFrom(int64(nd.Id)),
				IsPrivate: null.BoolFrom(false),
			},
		},
		{
			Name:               "Add Dir Document Fail, Name Validation Fail",
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			In: document.AddDirDocumentIn{
				DirId:     null.IntFrom(int64(nd.Id)),
				IsPrivate: null.BoolFrom(false),
			},
		},
		{
			Name:               "Add Dir Document Fail, Dir Id Validation Fail",
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			In: document.AddDirDocumentIn{
				Name:      "Dir A",
				IsPrivate: null.BoolFrom(false),
			},
		},

		{
			Name:               "Add Dir Document Fail, Private Status Validation Fail",
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			In: document.AddDirDocumentIn{
				Name:  "Dir A",
				DirId: null.IntFrom(int64(nd.Id)),
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
			res := documentDeps.AddDirDocument(ctx, c.In)
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

func TestAddFileDocument(t *testing.T) {
	err := ClearTables(db)
	if err != nil {
		t.Fatal(err)
	}

	nd, err := documentRepository.Save(context.Background(), dirSeed)
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
		In                 document.AddFileDocumentIn
	}{
		{
			Name:               "Add File Document at root dir Success",
			ExpectedStatusCode: http.StatusCreated,
			In: document.AddFileDocumentIn{
				DirId: null.IntFrom(0),
				File: httpdecode.FileHeader{
					Filename: fileName,
					File:     f,
				},
				IsPrivate: null.BoolFrom(false),
			},
		},
		{
			Name:               "Add File Document at 'n' dir Success",
			ExpectedStatusCode: http.StatusCreated,
			In: document.AddFileDocumentIn{
				DirId: null.IntFrom(int64(nd.Id)),
				File: httpdecode.FileHeader{
					Filename: fileName,
					File:     f,
				},
				IsPrivate: null.BoolFrom(false),
			},
		},
		{
			Name:               "Add File Document Fail, Dir Id Validation Fail",
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			In: document.AddFileDocumentIn{
				File: httpdecode.FileHeader{
					Filename: fileName,
					File:     f,
				},
				IsPrivate: null.BoolFrom(false),
			},
		},
		{
			Name:               "Add File Document Fail, File Validation Fail",
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			In: document.AddFileDocumentIn{
				DirId:     null.IntFrom(int64(nd.Id)),
				IsPrivate: null.BoolFrom(false),
			},
		},
		{
			Name:               "Add File Document Fail, Private Status Validation Fail",
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			In: document.AddFileDocumentIn{
				File: httpdecode.FileHeader{
					Filename: fileName,
					File:     f,
				},
				DirId: null.IntFrom(int64(nd.Id)),
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
			res := documentDeps.AddFileDocument(ctx, c.In)
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

	_, err = documentRepository.Save(context.Background(), dirSeed)
	if err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		Name               string
		ExpectedStatusCode int
	}{
		{
			Name:               "Query Documents Success",
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
			res := documentDeps.QueryDocument(ctx, "", "0")
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

func TestEditDirDocument(t *testing.T) {
	err := ClearTables(db)
	if err != nil {
		t.Fatal(err)
	}

	nf, err := documentRepository.Save(context.Background(), fileSeed)
	if err != nil {
		t.Fatal(err)
	}

	nd, err := documentRepository.Save(context.Background(), dirSeed)
	if err != nil {
		t.Fatal(err)
	}

	nid := strconv.FormatUint(nd.Id, 10)
	fid := strconv.FormatUint(nf.Id, 10)

	testCases := []struct {
		Name               string
		ExpectedStatusCode int
		Id                 string
		In                 document.EditDirDocumentIn
	}{
		{
			Name:               "Edit Dir Document Success",
			ExpectedStatusCode: http.StatusOK,
			Id:                 nid,
			In: document.EditDirDocumentIn{
				Name:      "Dir A",
				IsPrivate: null.BoolFrom(false),
			},
		},
		{
			Name:               "Edit Dir Document Fail, Not Found",
			ExpectedStatusCode: http.StatusNotFound,
			Id:                 "999",
			In: document.EditDirDocumentIn{
				Name:      "Dir A",
				IsPrivate: null.BoolFrom(false),
			},
		},
		{
			Name:               "Edit Dir Document Fail, Id is not Dir Id",
			ExpectedStatusCode: http.StatusNotFound,
			Id:                 fid,
			In: document.EditDirDocumentIn{
				Name:      "Dir A",
				IsPrivate: null.BoolFrom(false),
			},
		},
		{
			Name:               "Add Dir Document Fail, Name Validation Fail",
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			Id:                 nid,
			In: document.EditDirDocumentIn{
				IsPrivate: null.BoolFrom(false),
			},
		},
		{
			Name:               "Add Document Fail, Private Status Validation Fail",
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			Id:                 nid,
			In: document.EditDirDocumentIn{
				Name: "Dir A",
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
			res := documentDeps.EditDirDocument(ctx, c.Id, c.In)
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

func TestEditFileDocument(t *testing.T) {
	err := ClearTables(db)
	if err != nil {
		t.Fatal(err)
	}

	f, err := os.OpenFile(fileDir, os.O_RDONLY, 0o444)
	if err != nil {
		t.Fatal(err)
	}

	nf, err := documentRepository.Save(context.Background(), fileSeed)
	if err != nil {
		t.Fatal(err)
	}

	nd, err := documentRepository.Save(context.Background(), dirSeed)
	if err != nil {
		t.Fatal(err)
	}

	nid := strconv.FormatUint(nd.Id, 10)
	fid := strconv.FormatUint(nf.Id, 10)

	testCases := []struct {
		Name               string
		ExpectedStatusCode int
		Id                 string
		In                 document.EditFileDocumentIn
	}{
		{
			Name:               "Edit File Document Success",
			ExpectedStatusCode: http.StatusOK,
			Id:                 fid,
			In: document.EditFileDocumentIn{
				File: httpdecode.FileHeader{
					Filename: fileName,
					File:     f,
				},
				IsPrivate: null.BoolFrom(false),
			},
		},
		{
			Name:               "Edit File Document Fail, Not Found",
			ExpectedStatusCode: http.StatusNotFound,
			Id:                 "999",
			In: document.EditFileDocumentIn{
				File: httpdecode.FileHeader{
					Filename: fileName,
					File:     f,
				},
				IsPrivate: null.BoolFrom(false),
			},
		},
		{
			Name:               "Edit File Document Fail, Id is not File id",
			ExpectedStatusCode: http.StatusNotFound,
			Id:                 nid,
			In: document.EditFileDocumentIn{
				File: httpdecode.FileHeader{
					Filename: fileName,
					File:     f,
				},
				IsPrivate: null.BoolFrom(false),
			},
		},
		{
			Name:               "Edit File Document Fail, File Validation Fail",
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			Id:                 fid,
			In: document.EditFileDocumentIn{
				IsPrivate: null.BoolFrom(false),
			},
		},
		{
			Name:               "Edit File Document Fail, Private Status Validation Fail",
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			Id:                 fid,
			In: document.EditFileDocumentIn{
				File: httpdecode.FileHeader{
					Filename: fileName,
					File:     f,
				},
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
			res := documentDeps.EditFileDocument(ctx, c.Id, c.In)
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

func TestRemoveDocument(t *testing.T) {
	err := ClearTables(db)
	if err != nil {
		t.Fatal(err)
	}

	nf, err := documentRepository.Save(context.Background(), fileSeed)
	if err != nil {
		t.Fatal(err)
	}

	nd, err := documentRepository.Save(context.Background(), dirSeed)
	if err != nil {
		t.Fatal(err)
	}

	nid := strconv.FormatUint(nd.Id, 10)
	fid := strconv.FormatUint(nf.Id, 10)

	testCases := []struct {
		Name               string
		ExpectedStatusCode int
		Id                 string
	}{
		{
			Name:               "Remove Document (Dir) Success",
			ExpectedStatusCode: http.StatusOK,
			Id:                 nid,
		},
		{
			Name:               "Remove Document (File) Success",
			ExpectedStatusCode: http.StatusOK,
			Id:                 fid,
		},
		{
			Name:               "Remove Document (File/Dir) Fail, Not Found",
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
			res := documentDeps.RemoveDocument(ctx, c.Id)
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

func TestFindDocumentChildren(t *testing.T) {
	err := ClearTables(db)
	if err != nil {
		t.Fatal(err)
	}

	p, c, err := createDocumentChildren(documentRepository, dirSeed)
	if err != nil {
		t.Fatal(err)
	}

	pid := strconv.FormatUint(p.Id, 10)
	cid := strconv.FormatUint(c.Id, 10)

	testCases := []struct {
		Name               string
		ExpectedStatusCode int
		Id                 string
	}{
		{
			Name:               "Find Document (Dir) Childrens, Success",
			ExpectedStatusCode: http.StatusOK,
			Id:                 pid,
		},
		{
			Name:               "Find Document (Dir) Childrens, Success",
			ExpectedStatusCode: http.StatusOK,
			Id:                 cid,
		},
		{
			Name:               "Find Document (Dir) Childrens, Success",
			ExpectedStatusCode: http.StatusOK,
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
			res := documentDeps.FindDocumentChildren(ctx, c.Id, "")
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
