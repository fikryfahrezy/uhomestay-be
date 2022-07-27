package dues_test

import (
	"context"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"

	arbitary "github.com/PA-D3RPLA/d3if43-htt-uhomestay/arbitrary"
	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/dues"
	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/httpdecode"
	"gopkg.in/guregu/null.v4"
)

func TestQueryMemberDues(t *testing.T) {
	err := ClearTables(db)
	if err != nil {
		t.Fatal(err)
	}

	uid, _, err := createMemberNDues(duesDeps, memberSeed, duesSeed2)
	if err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		Name               string
		ExpectedStatusCode int
		Id                 string
	}{
		{
			Name:               "Query Member Dues Success",
			ExpectedStatusCode: http.StatusOK,
			Id:                 uid,
		},
		{
			Name:               "Query Member Dues Fail, Not Found",
			ExpectedStatusCode: http.StatusNotFound,
			Id:                 "12345678-1234-1234-1234-123456789012",
		},
		{
			Name:               "Query Member Dues Fail, Id not UUID",
			ExpectedStatusCode: http.StatusNotFound,
			Id:                 "blablabla",
		},
	}

	for _, c := range testCases {
		t.Run(c.Name, func(t *testing.T) {
			tx, err := db.Begin(context.Background())
			if err != nil {
				t.Fatal(err)
			}

			ctx := context.WithValue(context.Background(), arbitary.TrxX{}, tx)
			res := duesDeps.QueryMemberDues(ctx, c.Id, "", "")
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

func TestQueryMembersDues(t *testing.T) {
	err := ClearTables(db)
	if err != nil {
		t.Fatal(err)
	}

	nd1, err := duesRepository.Save(context.Background(), duesSeed)
	if err != nil {
		t.Fatal(err)
	}

	_, nd2, err := createMemberNDues(duesDeps, memberSeed, duesSeed2)
	if err != nil {
		t.Fatal(err)
	}

	pid1 := strconv.FormatUint(nd1.Id, 10)
	pid2 := strconv.FormatUint(nd2, 10)

	testCases := []struct {
		Name               string
		ExpectedStatusCode int
		Id                 string
	}{
		{
			Name:               "Query Members Dues Without Member in it Success",
			ExpectedStatusCode: http.StatusOK,
			Id:                 pid1,
		},
		{
			Name:               "Query Members Dues With Member in it Success",
			ExpectedStatusCode: http.StatusOK,
			Id:                 pid2,
		},
		{
			Name:               "Query Members Dues Success, Fallback to last Dues",
			ExpectedStatusCode: http.StatusOK,
			Id:                 "0",
		},
		{
			Name:               "Query Members Dues Success, Resulting Empty Data",
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
			res := duesDeps.QueryMembersDues(ctx, c.Id, dues.QueryMembersDuesQIn{})
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

func TestPayMemberDues(t *testing.T) {
	err := ClearTables(db)
	if err != nil {
		t.Fatal(err)
	}

	f, err := os.OpenFile(fileDir, os.O_RDONLY, 0o444)
	if err != nil {
		t.Fatal(err)
	}

	uid, _, nd2, err := createMemberDues(duesDeps, memberSeed, duesSeed2, unpaidMemDSeed)
	if err != nil {
		t.Fatal(err)
	}

	uid2, _, nd3, err := createMemberDues(duesDeps, memberSeed2, duesSeed3, paidMemDSeed)
	if err != nil {
		t.Fatal(err)
	}

	pid1 := strconv.FormatUint(nd2, 10)
	pid2 := strconv.FormatUint(nd3, 10)

	testCases := []struct {
		Name               string
		ExpectedStatusCode int
		Id                 string
		Uid                string
		In                 dues.PayMemberDuesIn
	}{
		{
			Name:               "Pay Member Dues Success",
			ExpectedStatusCode: http.StatusOK,
			Id:                 pid1,
			Uid:                uid,
			In: dues.PayMemberDuesIn{
				File: httpdecode.FileHeader{
					Filename: fileName,
					File:     f,
				},
			},
		},
		{
			Name:               "Pay Member Dues Fail, File Validation Fail",
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			Id:                 pid1,
			Uid:                uid,
			In:                 dues.PayMemberDuesIn{},
		},
		{
			Name:               "Pay Member Dues Fail, Dues Already Paid",
			ExpectedStatusCode: http.StatusNotFound,
			Id:                 pid2,
			Uid:                uid2,
			In: dues.PayMemberDuesIn{
				File: httpdecode.FileHeader{
					Filename: fileName,
					File:     f,
				},
			},
		},
		{
			Name:               "Pay Member Dues Fail, Dues Not Found",
			ExpectedStatusCode: http.StatusNotFound,
			Id:                 "999",
			Uid:                "12345678-1234-1234-1234-123456789012",
			In: dues.PayMemberDuesIn{
				File: httpdecode.FileHeader{
					Filename: fileName,
					File:     f,
				},
			},
		},
		{
			Name:               "Pay Member Dues with filename 200 chars, success",
			ExpectedStatusCode: http.StatusNotFound,
			Id:                 "999",
			Uid:                uid,
			In: dues.PayMemberDuesIn{
				File: httpdecode.FileHeader{
					Filename: strings.Repeat("a", 200),
					File:     f,
				},
			},
		},
		{
			Name:               "Pay Member Dues with filename over 200 chars, fail",
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			Id:                 "999",
			Uid:                uid,
			In: dues.PayMemberDuesIn{
				File: httpdecode.FileHeader{
					Filename: strings.Repeat("a", 201),
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
			res := duesDeps.PayMemberDues(ctx, c.Uid, c.Id, c.In)
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

func TestEditMemberDues(t *testing.T) {
	err := ClearTables(db)
	if err != nil {
		t.Fatal(err)
	}

	f, err := os.OpenFile(fileDir, os.O_RDONLY, 0o444)
	if err != nil {
		t.Fatal(err)
	}

	_, _, nd1, err := createMemberDues(duesDeps, memberSeed, duesSeed, unpaidMemDSeed)
	if err != nil {
		t.Fatal(err)
	}

	_, _, nd2, err := createMemberDues(duesDeps, memberSeed2, duesSeed2, paidMemDSeed)
	if err != nil {
		t.Fatal(err)
	}

	pid1 := strconv.FormatUint(nd1, 10)
	pid2 := strconv.FormatUint(nd2, 10)

	testCases := []struct {
		Name               string
		ExpectedStatusCode int
		Id                 string
		In                 dues.EditMemberDuesIn
	}{
		{
			Name:               "Edit Dues Success",
			ExpectedStatusCode: http.StatusOK,
			Id:                 pid1,
			In: dues.EditMemberDuesIn{
				File: httpdecode.FileHeader{
					Filename: fileName,
					File:     f,
				},
			},
		},
		{
			Name:               "Edit Dues Fail, Dues Already Paid",
			ExpectedStatusCode: http.StatusNotFound,
			Id:                 pid2,
			In: dues.EditMemberDuesIn{
				File: httpdecode.FileHeader{
					Filename: fileName,
					File:     f,
				},
			},
		},
		{
			Name:               "Edit Member Dues Fail, File Validation Fail",
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			Id:                 pid1,
			In:                 dues.EditMemberDuesIn{},
		},
		{
			Name:               "Edit Member Dues Fail, Not Found",
			ExpectedStatusCode: http.StatusNotFound,
			Id:                 "999",
			In: dues.EditMemberDuesIn{
				File: httpdecode.FileHeader{
					Filename: fileName,
					File:     f,
				},
			},
		},
		{
			Name:               "Edit Dues  with Filename 200 chars Success",
			ExpectedStatusCode: http.StatusOK,
			Id:                 pid1,
			In: dues.EditMemberDuesIn{
				File: httpdecode.FileHeader{
					Filename: strings.Repeat("a", 200),
					File:     f,
				},
			},
		},
		{
			Name:               "Edit Dues  with Filename over 200 chars fail",
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			Id:                 pid1,
			In: dues.EditMemberDuesIn{
				File: httpdecode.FileHeader{
					Filename: strings.Repeat("a", 201),
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
			res := duesDeps.EditMemberDues(ctx, c.Id, c.In)
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

func TestPaidMemberDues(t *testing.T) {
	err := ClearTables(db)
	if err != nil {
		t.Fatal(err)
	}

	_, _, nd1, err := createMemberDues(duesDeps, memberSeed, duesSeed, unpaidMemDSeed)
	if err != nil {
		t.Fatal(err)
	}

	_, _, nd2, err := createMemberDues(duesDeps, memberSeed2, duesSeed2, paidMemDSeed)
	if err != nil {
		t.Fatal(err)
	}

	pid1 := strconv.FormatUint(nd1, 10)
	pid2 := strconv.FormatUint(nd2, 10)

	testCases := []struct {
		Name               string
		ExpectedStatusCode int
		Id                 string
		In                 dues.PaidMemberDuesIn
	}{
		{
			Name:               "Paid Unpaid Member Dues to Paid Success",
			ExpectedStatusCode: http.StatusOK,
			Id:                 pid1,
			In: dues.PaidMemberDuesIn{
				IsPaid: null.BoolFrom(true),
			},
		},
		{
			Name:               "Paid Unpaid Member Dues to Paid Fail, Dues Already Paid",
			ExpectedStatusCode: http.StatusNotFound,
			Id:                 pid2,
			In: dues.PaidMemberDuesIn{
				IsPaid: null.BoolFrom(true),
			},
		},
		{
			Name:               "Paid Member Dues Fail, Is Paid Validation Fail",
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			Id:                 pid1,
			In:                 dues.PaidMemberDuesIn{},
		},
		{
			Name:               "Paid Member Dues Fail, Not Found",
			ExpectedStatusCode: http.StatusNotFound,
			Id:                 "999",
			In: dues.PaidMemberDuesIn{
				IsPaid: null.BoolFrom(true),
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
			res := duesDeps.PaidMemberDues(ctx, c.Id, c.In)
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
