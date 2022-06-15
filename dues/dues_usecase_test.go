package dues_test

import (
	"context"
	"net/http"
	"strconv"
	"testing"
	"time"

	arbitary "github.com/PA-D3RPLA/d3if43-htt-uhomestay/arbitrary"
	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/dues"
)

func TestAddDues(t *testing.T) {
	err := ClearTables(db)
	if err != nil {
		t.Fatal(err)
	}

	_, err = duesRepository.Save(context.Background(), duesSeed)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(pastDuesSeed.Date.Format("2006-01-02"))
	testCases := []struct {
		Name               string
		ExpectedStatusCode int
		In                 dues.AddDuesIn
	}{
		{
			Name:               "Add Dues Success",
			ExpectedStatusCode: http.StatusCreated,
			In: dues.AddDuesIn{
				Date:      time.Now().Format("2006-01-02"),
				IdrAmount: "100000",
			},
		},
		{
			Name:               "Add Dues Fail, Past Date",
			ExpectedStatusCode: http.StatusBadRequest,
			In: dues.AddDuesIn{
				Date:      pastDuesSeed.Date.Format("2006-01-02"),
				IdrAmount: "100000",
			},
		},
		{
			Name:               "Add Dues Fail, Date Already Occupied",
			ExpectedStatusCode: http.StatusBadRequest,
			In: dues.AddDuesIn{
				Date:      duesSeed.Date.Format("2006-01-02"),
				IdrAmount: "100000",
			},
		},
		{
			Name:               "Add Dues Fail, Date Validation Fail",
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			In: dues.AddDuesIn{
				IdrAmount: "100000",
			},
		},
		{
			Name:               "Add Dues Fail, Idr Amount Validation Fail",
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			In: dues.AddDuesIn{
				Date: time.Now().Format("2006-01-02"),
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
			res := duesDeps.AddDues(ctx, c.In)
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

func TestQueryDues(t *testing.T) {
	err := ClearTables(db)
	if err != nil {
		t.Fatal(err)
	}

	_, err = duesRepository.Save(context.Background(), duesSeed)
	if err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		Name               string
		ExpectedStatusCode int
	}{
		{
			Name:               "Query Dues Success",
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
			res := duesDeps.QueryDues(ctx, "")
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

func TestEditDues(t *testing.T) {
	err := ClearTables(db)
	if err != nil {
		t.Fatal(err)
	}

	nd, err := duesRepository.Save(context.Background(), duesSeed)
	if err != nil {
		t.Fatal(err)
	}

	_, nd2id, _, err := createMemberDues(duesDeps, memberSeed, duesSeed2, paidMemDSeed)
	if err != nil {
		t.Fatal(err)
	}

	pid := strconv.FormatUint(nd.Id, 10)
	pid2 := strconv.FormatUint(nd2id, 10)

	testCases := []struct {
		Name               string
		ExpectedStatusCode int
		Id                 string
		In                 dues.EditDuesIn
	}{
		{
			Name:               "Edit Dues Success",
			ExpectedStatusCode: http.StatusOK,
			Id:                 pid,
			In: dues.EditDuesIn{
				Date:      time.Now().Format("2006-01-02"),
				IdrAmount: "100000",
			},
		},
		{
			Name:               "Edit Dues Success, Change Nominal",
			ExpectedStatusCode: http.StatusOK,
			Id:                 pid,
			In: dues.EditDuesIn{
				Date:      time.Now().Format("2006-01-02"),
				IdrAmount: "200000",
			},
		},
		{
			Name:               "Edit Dues Fail, Date Already Occupied",
			ExpectedStatusCode: http.StatusBadRequest,
			Id:                 pid,
			In: dues.EditDuesIn{
				Date:      duesSeed2.Date.Format("2006-01-02"),
				IdrAmount: "100000",
			},
		},
		{
			Name:               "Edit Dues Fail, Somone Already Paid",
			ExpectedStatusCode: http.StatusBadRequest,
			Id:                 pid2,
			In: dues.EditDuesIn{
				Date:      duesSeed2.Date.Format("2006-01-02"),
				IdrAmount: "100000",
			},
		},
		{
			Name:               "Edit Dues Fail, Date Validation Fail",
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			Id:                 pid,
			In: dues.EditDuesIn{
				IdrAmount: "100000",
			},
		},
		{
			Name:               "Edit Dues Fail, Idr Amount Validation Fail",
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			Id:                 pid,
			In: dues.EditDuesIn{
				Date: time.Now().Format("2006-01-02"),
			},
		},
		{
			Name:               "Edit Dues Fail, Not Found",
			ExpectedStatusCode: http.StatusNotFound,
			Id:                 "999",
			In: dues.EditDuesIn{
				Date:      time.Now().Format("2006-01-02"),
				IdrAmount: "100000",
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
			res := duesDeps.EditDues(ctx, c.Id, c.In)
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

func TestRemoveDues(t *testing.T) {
	err := ClearTables(db)
	if err != nil {
		t.Fatal(err)
	}

	nd, err := duesRepository.Save(context.Background(), duesSeed)
	if err != nil {
		t.Fatal(err)
	}

	_, nd2id, _, err := createMemberDues(duesDeps, memberSeed, duesSeed2, paidMemDSeed)
	if err != nil {
		t.Fatal(err)
	}

	pid := strconv.FormatUint(nd.Id, 10)
	pid2 := strconv.FormatUint(nd2id, 10)

	testCases := []struct {
		Name               string
		ExpectedStatusCode int
		Id                 string
	}{
		{
			Name:               "Remove Dues Success",
			ExpectedStatusCode: http.StatusOK,
			Id:                 pid,
		},
		{
			Name:               "Remove Dues Fail, Someone Already Paid",
			ExpectedStatusCode: http.StatusBadRequest,
			Id:                 pid2,
		},
		{
			Name:               "Remove Dues Fail, Not Found",
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
			res := duesDeps.RemoveDues(ctx, c.Id)
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

func TestCheckPaidDues(t *testing.T) {
	err := ClearTables(db)
	if err != nil {
		t.Fatal(err)
	}

	_, nd1, _, err := createMemberDues(duesDeps, memberSeed, duesSeed, unpaidMemDSeed)
	if err != nil {
		t.Fatal(err)
	}

	_, nd2, _, err := createMemberDues(duesDeps, memberSeed2, duesSeed2, paidMemDSeed)
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
			Name:               "Check Paid Dues Success",
			ExpectedStatusCode: http.StatusOK,
			Id:                 pid1,
		},
		{
			Name:               "Check Paid Dues Success",
			ExpectedStatusCode: http.StatusOK,
			Id:                 pid2,
		},
		{
			Name:               "Check Paid Dues Fail, Not Found",
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
			res := duesDeps.CheckPaidDues(ctx, c.Id)
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
