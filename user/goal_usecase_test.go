package user_test

import (
	"context"
	"net/http"
	"strconv"
	"testing"

	arbitary "github.com/PA-D3RPLA/d3if43-htt-uhomestay/arbitrary"
	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/user"
)

func TestAddGoal(t *testing.T) {
	err := ClearTables(db)
	if err != nil {
		t.Fatal(err)
	}

	pr, err := orgPeriodRepository.Save(context.Background(), period)
	if err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		Name               string
		ExpectedStatusCode int
		In                 user.AddGoalIn
	}{
		{
			Name:               "Add Goal Success",
			ExpectedStatusCode: http.StatusCreated,
			In: user.AddGoalIn{
				Vision:      `{"test": "test"}`,
				VisionText:  "test",
				Mission:     `{"test": "test"}`,
				MissionText: "test",
				OrgPeriodId: int64(pr.Id),
			},
		},
		{
			Name:               "Add Goal Fail, Org Period Id Validation Fail",
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			In: user.AddGoalIn{
				Vision:      `{"test": "test"}`,
				VisionText:  "test",
				Mission:     `{"test": "test"}`,
				MissionText: "test",
			},
		},
		{
			Name:               "Add Goal Fail, Org Period Id Not Found",
			ExpectedStatusCode: http.StatusNotFound,
			In: user.AddGoalIn{
				Vision:      `{"test": "test"}`,
				VisionText:  "test",
				Mission:     `{"test": "test"}`,
				MissionText: "test",
				OrgPeriodId: 999,
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
			res := userDeps.AddGoal(ctx, c.In)
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

func TestFindOrgPeriodGoal(t *testing.T) {
	err := ClearTables(db)
	if err != nil {
		t.Fatal(err)
	}

	pr, err := orgPeriodRepository.Save(context.Background(), period)
	if err != nil {
		t.Fatal(err)
	}

	gscp := user.GoalModel(goalSeed)
	gscp.OrgPeriodId = pr.Id

	_, err = goalRepository.Save(context.Background(), gscp)
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
			Name:               "Find Goal Success",
			ExpectedStatusCode: http.StatusOK,
			Id:                 pid,
		},
		{
			Name:               "Find Goal Success, Fallback to Current Active Period",
			ExpectedStatusCode: http.StatusOK,
			Id:                 "0",
		},
		{
			Name:               "Find Goal without Org Period Id Success (Empty Content)",
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
			res := userDeps.FindOrgPeriodGoal(ctx, c.Id)
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
