package user_test

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"testing"

	arbitary "github.com/PA-D3RPLA/d3if43-htt-uhomestay/arbitrary"
	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/user"
)

func TestAddPosition(t *testing.T) {
	err := ClearTables(db)
	if err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		Name               string
		ExpectedStatusCode int
		In                 user.AddPositionIn
	}{
		{
			Name:               "Add Position Success",
			ExpectedStatusCode: http.StatusCreated,
			In: user.AddPositionIn{
				Name:  "Leader",
				Level: 1,
			},
		},
		{
			Name:               "Add Position Fail, Name Validaton Fail",
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			In: user.AddPositionIn{
				Level: 1,
			},
		},
		{
			Name:               "Add Position Fail, Level Validaton Fail",
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			In: user.AddPositionIn{
				Name: "Leader",
			},
		},
		{
			Name:               "Add Position with name 200 chars Success",
			ExpectedStatusCode: http.StatusCreated,
			In: user.AddPositionIn{
				Name:  strings.Repeat("a", 200),
				Level: 1,
			},
		},
		{
			Name:               "Add Position with name over 200 chars Success",
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			In: user.AddPositionIn{
				Name:  strings.Repeat("a", 201),
				Level: 1,
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
			res := userDeps.AddPosition(ctx, c.In)
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

func TestQueryPosition(t *testing.T) {
	err := ClearTables(db)
	if err != nil {
		t.Fatal(err)
	}

	_, err = positionRepository.Save(context.Background(), position)
	if err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		Name               string
		ExpectedStatusCode int
	}{
		{
			Name:               "Query Position Success",
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
			res := userDeps.QueryPosition(ctx, "", "0")
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

func TestQueryPositionLevel(t *testing.T) {
	err := ClearTables(db)
	if err != nil {
		t.Fatal(err)
	}

	_, err = positionRepository.Save(context.Background(), position)
	if err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		Name               string
		ExpectedStatusCode int
	}{
		{
			Name:               "Query Position Level Success",
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
			res := userDeps.QueryPosition(ctx, "", "0")
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

func TestEditPosition(t *testing.T) {
	err := ClearTables(db)
	if err != nil {
		t.Fatal(err)
	}

	ps, err := positionRepository.Save(context.Background(), position)
	if err != nil {
		t.Fatal(err)
	}

	pid := strconv.FormatUint(ps.Id, 10)
	testCases := []struct {
		Name               string
		ExpectedStatusCode int
		Id                 string
		In                 user.EditPositionIn
	}{
		{
			Name:               "Edit Position Success",
			ExpectedStatusCode: http.StatusOK,
			Id:                 pid,
			In: user.EditPositionIn{
				Name:  "Leader",
				Level: 1,
			},
		},
		{
			Name:               "Edit Position Fail, Position Not Found",
			ExpectedStatusCode: http.StatusNotFound,
			Id:                 "999",
			In: user.EditPositionIn{
				Name:  "Leader",
				Level: 1,
			},
		},
		{
			Name:               "Edit Position Fail, Name Validaton Fail",
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			Id:                 pid,
			In: user.EditPositionIn{
				Level: 1,
			},
		},
		{
			Name:               "Edit Position Fail, Level Validaton Fail",
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			Id:                 pid,
			In: user.EditPositionIn{
				Name: "Leader",
			},
		},
		{
			Name:               "Edit Position with name 200 chars Success",
			ExpectedStatusCode: http.StatusOK,
			Id:                 pid,
			In: user.EditPositionIn{
				Name:  strings.Repeat("a", 200),
				Level: 1,
			},
		},
		{
			Name:               "Edit Position with name over 200 chars fail",
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			Id:                 pid,
			In: user.EditPositionIn{
				Name:  strings.Repeat("a", 201),
				Level: 1,
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
			res := userDeps.EditPosition(ctx, c.Id, c.In)
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

func TestRemovePosition(t *testing.T) {
	err := ClearTables(db)
	if err != nil {
		t.Fatal(err)
	}

	ps, err := positionRepository.Save(context.Background(), position)
	if err != nil {
		t.Fatal(err)
	}

	pid := strconv.FormatUint(ps.Id, 10)
	testCases := []struct {
		Name               string
		ExpectedStatusCode int
		Id                 string
	}{
		{
			Name:               "Remove Position Success",
			ExpectedStatusCode: http.StatusOK,
			Id:                 pid,
		},
		{
			Name:               "Remove Position Fail, Position Not Found",
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
			res := userDeps.RemovePosition(ctx, c.Id)
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
