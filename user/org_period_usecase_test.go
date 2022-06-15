package user_test

import (
	"context"
	"net/http"
	"strconv"
	"testing"
	"time"

	arbitary "github.com/PA-D3RPLA/d3if43-htt-uhomestay/arbitrary"
	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/user"
	"gopkg.in/guregu/null.v4"
)

func TestAddPeriod(t *testing.T) {
	err := ClearTables(db)
	if err != nil {
		t.Fatal(err)
	}

	uid, err := createUser(memberRepository, member)
	if err != nil {
		t.Fatal(err)
	}

	ps, err := positionRepository.Save(context.Background(), position)
	if err != nil {
		t.Fatal(err)
	}

	startDate := time.Now().Format("2006-01-02")
	endDate := time.Now().Add(time.Hour * 24 * 365).Format("2006-01-02")

	testCases := []struct {
		Name               string
		ExpectedStatusCode int
		In                 user.AddPeriodIn
	}{
		{
			Name:               "Add Period Susccess",
			ExpectedStatusCode: http.StatusCreated,
			In: user.AddPeriodIn{
				StartDate: startDate,
				EndDate:   endDate,
				Positions: []user.PositionIn{
					{
						Id: int64(ps.Id),
						Members: []user.MemberIn{
							{
								Id: uid,
							},
						},
					},
				},
			},
		},
		{
			Name:               "Add Period Fail, Start Date Validation Fail",
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			In: user.AddPeriodIn{
				StartDate: "",
				EndDate:   endDate,
				Positions: []user.PositionIn{},
			},
		},
		{
			Name:               "Add Period Fail, End Date Validation Fail",
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			In: user.AddPeriodIn{
				StartDate: startDate,
				EndDate:   "",
				Positions: []user.PositionIn{},
			},
		},
		{
			Name:               "Add Period Fail, Start Date Not Date",
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			In: user.AddPeriodIn{
				StartDate: "blabla",
				EndDate:   endDate,
				Positions: []user.PositionIn{},
			},
		},
		{
			Name:               "Add Period Fail, End Date Not Date",
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			In: user.AddPeriodIn{
				StartDate: startDate,
				EndDate:   "blabla",
				Positions: []user.PositionIn{},
			},
		},
		{
			Name:               "Add Period Fail, End Date is Lower Than Start Date",
			ExpectedStatusCode: http.StatusBadRequest,
			In: user.AddPeriodIn{
				StartDate: endDate,
				EndDate:   startDate,
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
			res := userDeps.AddPeriod(ctx, c.In)
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

func TestQueryPeriod(t *testing.T) {
	err := ClearTables(db)
	if err != nil {
		t.Fatal(err)
	}

	_, err = createOrgStructute(userDeps, member, period, position)
	if err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		Name               string
		ExpectedStatusCode int
	}{
		{
			Name:               "Query Period Susccess",
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
			res := userDeps.QueryPeriod(ctx, "")
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

func TestEditPeriod(t *testing.T) {
	err := ClearTables(db)
	if err != nil {
		t.Fatal(err)
	}

	uid, err := createUser(memberRepository, member2)
	if err != nil {
		t.Fatal(err)
	}

	ps, err := positionRepository.Save(context.Background(), position)
	if err != nil {
		t.Fatal(err)
	}

	ors, err := createOrgStructute(userDeps, member, period, position)
	if err != nil {
		t.Fatal(err)
	}

	ors2, err := createOrgStructute(userDeps, memberNormal, inactivePeriod, position)
	if err != nil {
		t.Fatal(err)
	}

	startDate := time.Now().Format("2006-01-02")
	endDate := time.Now().Add(time.Hour * 24 * 365).Format("2006-01-02")

	pid := strconv.FormatUint(ors.Id, 10)
	pid2 := strconv.FormatUint(ors2.Id, 10)

	testCases := []struct {
		Name               string
		ExpectedStatusCode int
		Id                 string
		In                 user.EditPeriodIn
	}{
		{
			Name:               "Edit Period Susccess",
			ExpectedStatusCode: http.StatusOK,
			Id:                 pid,
			In: user.EditPeriodIn{
				StartDate: startDate,
				EndDate:   endDate,
				Positions: []user.PositionIn{
					{
						Id: int64(ps.Id),
						Members: []user.MemberIn{
							{
								Id: uid,
							},
						},
					},
				},
			},
		},
		{
			Name:               "Edit Period Fail, Period not Found",
			ExpectedStatusCode: http.StatusNotFound,
			Id:                 "999",
			In: user.EditPeriodIn{
				StartDate: startDate,
				EndDate:   endDate,
				Positions: []user.PositionIn{},
			},
		},
		{
			Name:               "Edit Period Fail, Period not Found because not Active",
			ExpectedStatusCode: http.StatusNotFound,
			Id:                 pid2,
			In: user.EditPeriodIn{
				StartDate: startDate,
				EndDate:   endDate,
				Positions: []user.PositionIn{},
			},
		},
		{
			Name:               "Edit Period Fail, Start Date Validation Fail",
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			Id:                 pid,
			In: user.EditPeriodIn{
				StartDate: "",
				EndDate:   endDate,
				Positions: []user.PositionIn{},
			},
		},
		{
			Name:               "Edit Period Fail, End Date Validation Fail",
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			Id:                 pid,
			In: user.EditPeriodIn{
				StartDate: startDate,
				EndDate:   "",
				Positions: []user.PositionIn{},
			},
		},
		{
			Name:               "Edit Period Fail, Start Date Not Date",
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			Id:                 pid,
			In: user.EditPeriodIn{
				StartDate: "blabla",
				EndDate:   endDate,
				Positions: []user.PositionIn{},
			},
		},
		{
			Name:               "Edit Period Fail, End Date Not Date",
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			Id:                 pid,
			In: user.EditPeriodIn{
				StartDate: startDate,
				EndDate:   "blabla",
				Positions: []user.PositionIn{},
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
			res := userDeps.EditPeriod(ctx, c.Id, c.In)
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

func TestRemovePeriod(t *testing.T) {
	err := ClearTables(db)
	if err != nil {
		t.Fatal(err)
	}

	pr, err := orgPeriodRepository.Save(context.Background(), period)
	if err != nil {
		t.Fatal(err)
	}

	pid := strconv.FormatUint(pr.Id, 10)
	testCases := []struct {
		Name               string
		ExpectedStatusCode int
		Id                 string
	}{
		{
			Name:               "Remove Period Susccess",
			ExpectedStatusCode: http.StatusOK,
			Id:                 pid,
		},
		{
			Name:               "Remove Period Fail, Period not Found",
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
			res := userDeps.RemovePeriod(ctx, c.Id)
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

func TestSwitchPeriodStatus(t *testing.T) {
	err := ClearTables(db)
	if err != nil {
		t.Fatal(err)
	}

	pr, err := orgPeriodRepository.Save(context.Background(), period)
	if err != nil {
		t.Fatal(err)
	}

	pid := strconv.FormatUint(pr.Id, 10)
	testCases := []struct {
		Name               string
		ExpectedStatusCode int
		Id                 string
		In                 user.SwitchPeriodStatusIn
	}{
		{
			Name:               "Swtich Period Status Susccess",
			ExpectedStatusCode: http.StatusOK,
			Id:                 pid,
			In: user.SwitchPeriodStatusIn{
				IsActive: null.BoolFrom(false),
			},
		},
		{
			Name:               "Swtich Period Status Fail, Period Not Found",
			ExpectedStatusCode: http.StatusNotFound,
			Id:                 "999",
			In: user.SwitchPeriodStatusIn{
				IsActive: null.BoolFrom(false),
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
			res := userDeps.SwitchPeriodStatus(ctx, c.Id, c.In)
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

func TestFindActivePeriod(t *testing.T) {
	err := ClearTables(db)
	if err != nil {
		t.Fatal(err)
	}

	_, err = orgPeriodRepository.Save(context.Background(), period)
	if err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		Name               string
		ExpectedStatusCode int
	}{
		{
			Name:               "Query Period Susccess",
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
			res := userDeps.FindActivePeriod(ctx)
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

func TestQueryPeriodStructure(t *testing.T) {
	err := ClearTables(db)
	if err != nil {
		t.Fatal(err)
	}

	ors, err := createOrgStructute(userDeps, member, period, position)
	if err != nil {
		t.Fatal(err)
	}

	pid := strconv.FormatUint(ors.Id, 10)
	testCases := []struct {
		Name               string
		ExpectedStatusCode int
		Id                 string
	}{
		{
			Name:               "Query Active Period Susccess",
			ExpectedStatusCode: http.StatusOK,
			Id:                 pid,
		},
		{
			Name:               "Query Period Susccess Fallback to Active Period",
			ExpectedStatusCode: http.StatusOK,
			Id:                 "0",
		},
		{
			Name:               "Query Active Period Success, Resulting Empty Data",
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
			res := userDeps.QueryPeriodStructure(ctx, c.Id)
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
