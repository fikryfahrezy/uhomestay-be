package cashflow_test

import (
	"context"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	arbitary "github.com/PA-D3RPLA/d3if43-htt-uhomestay/arbitrary"
	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/cashflow"
	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/httpdecode"
)

func TestAddCashflow(t *testing.T) {
	err := ClearTables(db)
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
		In                 cashflow.AddCashflowIn
	}{
		{
			Name:               "Add Income Cashflow Without Note and File Success",
			ExpectedStatusCode: http.StatusCreated,
			In: cashflow.AddCashflowIn{
				Date:      time.Now().Format("2006-01-02"),
				IdrAmount: "10000",
				Type:      "income",
			},
		},
		{
			Name:               "Add Income Cashflow Without Note Success",
			ExpectedStatusCode: http.StatusCreated,
			In: cashflow.AddCashflowIn{
				Date:      time.Now().Format("2006-01-02"),
				IdrAmount: "10000",
				Type:      "income",
				File: httpdecode.FileHeader{
					Filename: fileName,
					File:     f,
				},
			},
		},
		{
			Name:               "Add Income Cashflow Without File Success",
			ExpectedStatusCode: http.StatusCreated,
			In: cashflow.AddCashflowIn{
				Date:      time.Now().Format("2006-01-02"),
				IdrAmount: "10000",
				Type:      "income",
				Note:      "Just Note",
			},
		},
		{
			Name:               "Add Outcome Cashflow Without Note and File Success",
			ExpectedStatusCode: http.StatusCreated,
			In: cashflow.AddCashflowIn{
				Date:      time.Now().Format("2006-01-02"),
				IdrAmount: "10000",
				Type:      "outcome",
			},
		},
		{
			Name:               "Add Outcome Cashflow Without Note Success",
			ExpectedStatusCode: http.StatusCreated,
			In: cashflow.AddCashflowIn{
				Date:      time.Now().Format("2006-01-02"),
				IdrAmount: "10000",
				Type:      "outcome",
				File: httpdecode.FileHeader{
					Filename: fileName,
					File:     f,
				},
			},
		},
		{
			Name:               "Add Outcome Cashflow Without File Success",
			ExpectedStatusCode: http.StatusCreated,
			In: cashflow.AddCashflowIn{
				Date:      time.Now().Format("2006-01-02"),
				IdrAmount: "10000",
				Type:      "outcome",
				Note:      "Just Note",
			},
		},
		{
			Name:               "Add Cashflow Fail, Unknown Type",
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			In: cashflow.AddCashflowIn{
				Date:      time.Now().Format("2006-01-02"),
				IdrAmount: "10000",
				Type:      "blabla",
			},
		},
		{
			Name:               "Add Cashflow Fail, Type Validation Fail",
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			In: cashflow.AddCashflowIn{
				Date:      time.Now().Format("2006-01-02"),
				IdrAmount: "10000",
			},
		},
		{
			Name:               "Add Cashflow Fail, Date Validation Fail",
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			In: cashflow.AddCashflowIn{
				IdrAmount: "10000",
				Type:      "income",
			},
		},
		{
			Name:               "Add Cashflow Fail, Idr Amount Validation Fail",
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			In: cashflow.AddCashflowIn{
				Date: time.Now().Format("2006-01-02"),
				Type: "income",
			},
		},
		{
			Name:               "Add Income Cashflow with Idr Amout 200 chars Success",
			ExpectedStatusCode: http.StatusCreated,
			In: cashflow.AddCashflowIn{
				Date:      time.Now().Format("2006-01-02"),
				IdrAmount: strings.Repeat("1", 200),
				Type:      "income",
				Note:      "Just Note",
			},
		},
		{
			Name:               "Add Income Cashflow with Idr Amout over 200 chars fail",
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			In: cashflow.AddCashflowIn{
				Date:      time.Now().Format("2006-01-02"),
				IdrAmount: strings.Repeat("1", 201),
				Type:      "income",
				Note:      "Just Note",
			},
		},
	}

	for _, c := range testCases {
		t.Run(c.Name, func(t *testing.T) {
			tx, err := db.Begin(context.Background())

			ctx := context.WithValue(context.Background(), arbitary.TrxX{}, tx)
			res := cashflowDeps.AddCashflow(ctx, c.In)
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

func TestQueryCashflow(t *testing.T) {
	err := ClearTables(db)
	if err != nil {
		t.Fatal(err)
	}

	_, err = cashflowRepository.Save(context.Background(), cashflowSeed)
	if err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		Name               string
		ExpectedStatusCode int
	}{
		{
			Name:               "Query Cashflows Success",
			ExpectedStatusCode: http.StatusOK,
		},
	}

	for _, c := range testCases {
		t.Run(c.Name, func(t *testing.T) {
			tx, err := db.Begin(context.Background())

			ctx := context.WithValue(context.Background(), arbitary.TrxX{}, tx)
			res := cashflowDeps.QueryCashflow(ctx, "", "")
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

func TestEditCashflow(t *testing.T) {
	err := ClearTables(db)
	if err != nil {
		t.Fatal(err)
	}

	f, err := os.OpenFile(fileDir, os.O_RDONLY, 0o444)
	if err != nil {
		t.Fatal(err)
	}

	nc, err := cashflowRepository.Save(context.Background(), cashflowSeed)
	if err != nil {
		t.Fatal(err)
	}

	pid := strconv.FormatUint(nc.Id, 10)

	testCases := []struct {
		Name               string
		ExpectedStatusCode int
		Id                 string
		In                 cashflow.EditCashflowIn
	}{
		{
			Name:               "Edit Income Cashflow Without Note and File Success",
			ExpectedStatusCode: http.StatusOK,
			Id:                 pid,
			In: cashflow.EditCashflowIn{
				Date:      time.Now().Format("2006-01-02"),
				IdrAmount: "10000",
				Type:      "income",
			},
		},
		{
			Name:               "Edit Income Cashflow Without Note Success",
			ExpectedStatusCode: http.StatusOK,
			Id:                 pid,
			In: cashflow.EditCashflowIn{
				Date:      time.Now().Format("2006-01-02"),
				IdrAmount: "10000",
				Type:      "income",
				File: httpdecode.FileHeader{
					Filename: fileName,
					File:     f,
				},
			},
		},
		{
			Name:               "Edit Income Cashflow Without File Success",
			ExpectedStatusCode: http.StatusOK,
			Id:                 pid,
			In: cashflow.EditCashflowIn{
				Date:      time.Now().Format("2006-01-02"),
				IdrAmount: "10000",
				Type:      "income",
				Note:      "Just Note",
			},
		},
		{
			Name:               "Edit Outcome Cashflow Without Note and File Success",
			ExpectedStatusCode: http.StatusOK,
			Id:                 pid,
			In: cashflow.EditCashflowIn{
				Date:      time.Now().Format("2006-01-02"),
				IdrAmount: "10000",
				Type:      "outcome",
			},
		},
		{
			Name:               "Edit Outcome Cashflow Without Note Success",
			ExpectedStatusCode: http.StatusOK,
			Id:                 pid,
			In: cashflow.EditCashflowIn{
				Date:      time.Now().Format("2006-01-02"),
				IdrAmount: "10000",
				Type:      "outcome",
				File: httpdecode.FileHeader{
					Filename: fileName,
					File:     f,
				},
			},
		},
		{
			Name:               "Edit Outcome Cashflow Without File Success",
			ExpectedStatusCode: http.StatusOK,
			Id:                 pid,
			In: cashflow.EditCashflowIn{
				Date:      time.Now().Format("2006-01-02"),
				IdrAmount: "10000",
				Type:      "outcome",
				Note:      "Just Note",
			},
		},
		{
			Name:               "Edit Cashflow Fail, Unknown Type",
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			Id:                 pid,
			In: cashflow.EditCashflowIn{
				Date:      time.Now().Format("2006-01-02"),
				IdrAmount: "10000",
				Type:      "blabla",
			},
		},
		{
			Name:               "Edit Cashflow Fail, Type Validation Fail",
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			Id:                 pid,
			In: cashflow.EditCashflowIn{
				Date:      time.Now().Format("2006-01-02"),
				IdrAmount: "10000",
			},
		},
		{
			Name:               "Edit Cashflow Fail, Date Validation Fail",
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			Id:                 pid,
			In: cashflow.EditCashflowIn{
				IdrAmount: "10000",
				Type:      "income",
			},
		},
		{
			Name:               "Edit Cashflow Fail, Idr Amount Validation Fail",
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			Id:                 pid,
			In: cashflow.EditCashflowIn{
				Date: time.Now().Format("2006-01-02"),
				Type: "income",
			},
		},
		{
			Name:               "Edit Income Cashflow Fail, Cashflow Not Found",
			ExpectedStatusCode: http.StatusNotFound,
			Id:                 "999",
			In: cashflow.EditCashflowIn{
				Date:      time.Now().Format("2006-01-02"),
				IdrAmount: "10000",
				Type:      "income",
			},
		},
		{
			Name:               "Edit Income Cashflow with Idr Amount 200 chars success",
			ExpectedStatusCode: http.StatusOK,
			Id:                 pid,
			In: cashflow.EditCashflowIn{
				Date:      time.Now().Format("2006-01-02"),
				IdrAmount: strings.Repeat("1", 200),
				Type:      "income",
				Note:      "Just Note",
			},
		},
		{
			Name:               "Edit Income Cashflow with Idr Amount over 200 chars fail",
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			Id:                 pid,
			In: cashflow.EditCashflowIn{
				Date:      time.Now().Format("2006-01-02"),
				IdrAmount: strings.Repeat("1", 201),
				Type:      "income",
				Note:      "Just Note",
			},
		},
	}

	for _, c := range testCases {
		t.Run(c.Name, func(t *testing.T) {
			tx, err := db.Begin(context.Background())

			ctx := context.WithValue(context.Background(), arbitary.TrxX{}, tx)
			res := cashflowDeps.EditCashflow(ctx, c.Id, c.In)
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

func TestRemoveCashflow(t *testing.T) {
	err := ClearTables(db)
	if err != nil {
		t.Fatal(err)
	}

	nc, err := cashflowRepository.Save(context.Background(), cashflowSeed)
	if err != nil {
		t.Fatal(err)
	}

	pid := strconv.FormatUint(nc.Id, 10)
	testCases := []struct {
		Name               string
		ExpectedStatusCode int
		Id                 string
	}{
		{
			Name:               "Remove Cashflow Success",
			ExpectedStatusCode: http.StatusOK,
			Id:                 pid,
		},
		{
			Name:               "Remove Cashflow, Not Found",
			ExpectedStatusCode: http.StatusNotFound,
			Id:                 "999",
		},
	}

	for _, c := range testCases {
		t.Run(c.Name, func(t *testing.T) {
			tx, err := db.Begin(context.Background())

			ctx := context.WithValue(context.Background(), arbitary.TrxX{}, tx)
			res := cashflowDeps.RemoveCashflow(ctx, c.Id)
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

func TestCalculateCashflow(t *testing.T) {
	err := ClearTables(db)
	if err != nil {
		t.Fatal(err)
	}

	_, err = cashflowRepository.Save(context.Background(), cashflowSeed)
	if err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		Name               string
		ExpectedStatusCode int
	}{
		{
			Name:               "Calculate Cashflows Success",
			ExpectedStatusCode: http.StatusOK,
		},
	}

	for _, c := range testCases {
		t.Run(c.Name, func(t *testing.T) {
			tx, err := db.Begin(context.Background())

			ctx := context.WithValue(context.Background(), arbitary.TrxX{}, tx)
			res := cashflowDeps.CalculateCashflow(ctx)
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
