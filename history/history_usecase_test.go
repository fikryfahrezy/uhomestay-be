package history_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/history"
)

func TestAddHistory(t *testing.T) {
	err := ClearTables(db)
	if err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		Name               string
		ExpectedStatusCode int
		In                 history.AddHistoryIn
	}{
		{
			Name:               "Add History Success",
			ExpectedStatusCode: http.StatusCreated,
			In: history.AddHistoryIn{
				Content:     `{"test": "hi"}`,
				ContentText: "hi",
			},
		},
	}

	for _, c := range testCases {
		t.Run(c.Name, func(t *testing.T) {
			res := historyDeps.AddHistory(context.Background(), c.In)

			if res.StatusCode != c.ExpectedStatusCode {
				t.Logf("%#v", res)
				t.Log(err)
				t.Fatalf("Expected response code %d. Got %d\n", c.ExpectedStatusCode, res.StatusCode)
			}
		})
	}
}

func TestFindLatestHistory(t *testing.T) {
	err := ClearTables(db)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("Get Latest History Success (Empty Conten)", func(t *testing.T) {
		res := historyDeps.FindLatestHistory(context.Background())

		if res.StatusCode != http.StatusOK {
			t.Logf("%#v", res)
			t.Log(err)
			t.Fatalf("Expected response code %d. Got %d\n", http.StatusOK, res.StatusCode)
		}
	})

	t.Run("Get Latest History Success", func(t *testing.T) {
		_, err := historyRepository.Save(context.Background(), historySeed)
		if err != nil {
			t.Fatal(err)
		}

		res := historyDeps.FindLatestHistory(context.Background())

		if res.StatusCode != http.StatusOK {
			t.Logf("%#v", res)
			t.Log(err)
			t.Fatalf("Expected response code %d. Got %d\n", http.StatusOK, res.StatusCode)
		}
	})
}
