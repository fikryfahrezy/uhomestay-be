package homestay_test

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"testing"

	arbitary "github.com/PA-D3RPLA/d3if43-htt-uhomestay/arbitrary"
	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/homestay"
)

func TestAddMemberHomestay(t *testing.T) {
	err := ClearTables(db)
	if err != nil {
		t.Fatal(err)
	}

	muid, err := createUser(memberRepository, memberSeed)
	if err != nil {
		t.Fatal(err)
	}

	imageId, err := createHomestayImage(homestayImageRepository, fileSeed)
	if err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		ExpectedStatusCode int
		Name               string
		Uid                string
		In                 homestay.AddMemberHomestayIn
	}{
		{
			Name:               "Add Member Homestay Success",
			ExpectedStatusCode: http.StatusCreated,
			Uid:                muid,
			In: homestay.AddMemberHomestayIn{
				Name:      "Homestay Name",
				Address:   "Homestay Address",
				Latitude:  "120.12312312",
				Longitude: "90.1212321",
				ImageIds:  []int64{imageId},
			},
		},
		{
			Name:               "Add Member Homestay Fail, Name Required",
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			Uid:                muid,
			In: homestay.AddMemberHomestayIn{
				Name:      "",
				Address:   "Homestay Address",
				Latitude:  "120.12312312",
				Longitude: "90.1212321",
				ImageIds:  []int64{imageId},
			},
		},
		{
			Name:               "Add Member Homestay Fail, Name over 100 characters",
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			Uid:                muid,
			In: homestay.AddMemberHomestayIn{
				Name:      strings.Repeat("a", 101),
				Address:   "Homestay Address",
				Latitude:  "120.12312312",
				Longitude: "90.1212321",
				ImageIds:  []int64{imageId},
			},
		},
		{
			Name:               "Add Member Homestay Fail, Address Required",
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			Uid:                muid,
			In: homestay.AddMemberHomestayIn{
				Name:      "Homestay Name",
				Address:   "",
				Latitude:  "120.12312312",
				Longitude: "90.1212321",
				ImageIds:  []int64{imageId},
			},
		},
		{
			Name:               "Add Member Homestay Fail, Address over 200 characters",
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			Uid:                muid,
			In: homestay.AddMemberHomestayIn{
				Name:      "Homestay Name",
				Address:   strings.Repeat("a", 201),
				Latitude:  "120.12312312",
				Longitude: "90.1212321",
				ImageIds:  []int64{imageId},
			},
		},
		{
			Name:               "Add Member Homestay Fail, Latitude Required",
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			Uid:                muid,
			In: homestay.AddMemberHomestayIn{
				Name:      "Homestay Name",
				Address:   "Homestay Address",
				Latitude:  "",
				Longitude: "90.1212321",
				ImageIds:  []int64{imageId},
			},
		},
		{
			Name:               "Add Member Homestay Fail, Latitude over 50 characters",
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			Uid:                muid,
			In: homestay.AddMemberHomestayIn{
				Name:      "Homestay Name",
				Address:   "Homestay Address",
				Latitude:  strings.Repeat("0", 51),
				Longitude: "90.1212321",
				ImageIds:  []int64{imageId},
			},
		},
		{
			Name:               "Add Member Homestay Fail, Longitude Required",
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			Uid:                muid,
			In: homestay.AddMemberHomestayIn{
				Name:      "Homestay Name",
				Address:   "Homestay Address",
				Latitude:  "120.12312312",
				Longitude: "",
				ImageIds:  []int64{imageId},
			},
		},
		{
			Name:               "Add Member Homestay Fail, Longitude over 50 characters",
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			Uid:                muid,
			In: homestay.AddMemberHomestayIn{
				Name:      "Homestay Name",
				Address:   "Homestay Address",
				Latitude:  "120.12312312",
				Longitude: strings.Repeat("0", 51),
				ImageIds:  []int64{imageId},
			},
		},
		{
			Name:               "Add Member Homestay Fail, Images not valid",
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			Uid:                muid,
			In: homestay.AddMemberHomestayIn{
				Name:      "Homestay Name",
				Address:   "Homestay Address",
				Latitude:  "120.12312312",
				Longitude: strings.Repeat("0", 51),
				ImageIds:  []int64{99},
			},
		},
		{
			Name:               "Add Member Homestay Fail, Member uid not valid",
			ExpectedStatusCode: http.StatusNotFound,
			Uid:                "",
			In: homestay.AddMemberHomestayIn{
				Name:      "Homestay Name",
				Address:   "Homestay Address",
				Latitude:  "120.12312312",
				Longitude: "90.1212321",
				ImageIds:  []int64{imageId},
			},
		},
		{
			Name:               "Add Member Homestay Fail, Member not Found",
			ExpectedStatusCode: http.StatusNotFound,
			Uid:                "12345678-1234-1234-1234-123456789012",
			In: homestay.AddMemberHomestayIn{
				Name:      "Homestay Name",
				Address:   "Homestay Address",
				Latitude:  "120.12312312",
				Longitude: "90.1212321",
				ImageIds:  []int64{imageId},
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
			res := homestayDeps.AddMemberHomestay(ctx, c.Uid, c.In)
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

func TestEditMemberHomestay(t *testing.T) {
	err := ClearTables(db)
	if err != nil {
		t.Fatal(err)
	}

	muid, err := createUser(memberRepository, memberSeed)
	if err != nil {
		t.Fatal(err)
	}

	imageId, err := createHomestayImage(homestayImageRepository, fileSeed)
	if err != nil {
		t.Fatal(err)
	}

	memberHomeId, err := createMemberHomestay(memberHomestayRepository, muid, homestaySeed)
	if err != nil {
		t.Fatal(err)
	}

	mhId := strconv.FormatInt(memberHomeId, 10)

	testCases := []struct {
		ExpectedStatusCode int
		Name               string
		Uid                string
		Pid                string
		In                 homestay.EditMemberHomestayIn
	}{
		{
			Name:               "Edit Member Homestay Success",
			ExpectedStatusCode: http.StatusOK,
			Uid:                muid,
			Pid:                mhId,
			In: homestay.EditMemberHomestayIn{
				Name:      "Homestay Name",
				Address:   "Homestay Address",
				Latitude:  "120.12312312",
				Longitude: "90.1212321",
				ImageIds:  []int64{imageId},
			},
		},
		{
			Name:               "Edit Member Homestay Fail, Name Required",
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			Uid:                muid,
			Pid:                mhId,
			In: homestay.EditMemberHomestayIn{
				Name:      "",
				Address:   "Homestay Address",
				Latitude:  "120.12312312",
				Longitude: "90.1212321",
				ImageIds:  []int64{imageId},
			},
		},
		{
			Name:               "Edit Member Homestay Fail, Name over 100 characters",
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			Uid:                muid,
			Pid:                mhId,
			In: homestay.EditMemberHomestayIn{
				Name:      strings.Repeat("a", 101),
				Address:   "Homestay Address",
				Latitude:  "120.12312312",
				Longitude: "90.1212321",
				ImageIds:  []int64{imageId},
			},
		},
		{
			Name:               "Edit Member Homestay Fail, Address Required",
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			Uid:                muid,
			Pid:                mhId,
			In: homestay.EditMemberHomestayIn{
				Name:      "Homestay Name",
				Address:   "",
				Latitude:  "120.12312312",
				Longitude: "90.1212321",
				ImageIds:  []int64{imageId},
			},
		},
		{
			Name:               "Edit Member Homestay Fail, Address over 200 characters",
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			Uid:                muid,
			Pid:                mhId,
			In: homestay.EditMemberHomestayIn{
				Name:      "Homestay Name",
				Address:   strings.Repeat("a", 201),
				Latitude:  "120.12312312",
				Longitude: "90.1212321",
				ImageIds:  []int64{imageId},
			},
		},
		{
			Name:               "Edit Member Homestay Fail, Latitude Required",
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			Uid:                muid,
			Pid:                mhId,
			In: homestay.EditMemberHomestayIn{
				Name:      "Homestay Name",
				Address:   "Homestay Address",
				Latitude:  "",
				Longitude: "90.1212321",
				ImageIds:  []int64{imageId},
			},
		},
		{
			Name:               "Edit Member Homestay Fail, Latitude over 50 characters",
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			Uid:                muid,
			Pid:                mhId,
			In: homestay.EditMemberHomestayIn{
				Name:      "Homestay Name",
				Address:   "Homestay Address",
				Latitude:  strings.Repeat("0", 51),
				Longitude: "90.1212321",
				ImageIds:  []int64{imageId},
			},
		},
		{
			Name:               "Edit Member Homestay Fail, Longitude Required",
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			Uid:                muid,
			Pid:                mhId,
			In: homestay.EditMemberHomestayIn{
				Name:      "Homestay Name",
				Address:   "Homestay Address",
				Latitude:  "120.12312312",
				Longitude: "",
				ImageIds:  []int64{imageId},
			},
		},
		{
			Name:               "Edit Member Homestay Fail, Longitude over 50 characters",
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			Uid:                muid,
			Pid:                mhId,
			In: homestay.EditMemberHomestayIn{
				Name:      "Homestay Name",
				Address:   "Homestay Address",
				Latitude:  "120.12312312",
				Longitude: strings.Repeat("0", 51),
				ImageIds:  []int64{imageId},
			},
		},
		{
			Name:               "Edit Member Homestay Fail, Images not valid",
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			Uid:                muid,
			Pid:                mhId,
			In: homestay.EditMemberHomestayIn{
				Name:      "Homestay Name",
				Address:   "Homestay Address",
				Latitude:  "120.12312312",
				Longitude: strings.Repeat("0", 51),
				ImageIds:  []int64{99},
			},
		},
		{
			Name:               "Edit Member Homestay Fail, Member uid not valid",
			ExpectedStatusCode: http.StatusNotFound,
			Uid:                "",
			Pid:                mhId,
			In: homestay.EditMemberHomestayIn{
				Name:      "Homestay Name",
				Address:   "Homestay Address",
				Latitude:  "120.12312312",
				Longitude: "90.1212321",
				ImageIds:  []int64{imageId},
			},
		},
		{
			Name:               "Edit Member Homestay Fail, Member not Found",
			ExpectedStatusCode: http.StatusNotFound,
			Uid:                "12345678-1234-1234-1234-123456789012",
			Pid:                mhId,
			In: homestay.EditMemberHomestayIn{
				Name:      "Homestay Name",
				Address:   "Homestay Address",
				Latitude:  "120.12312312",
				Longitude: "90.1212321",
				ImageIds:  []int64{imageId},
			},
		},
		{
			Name:               "Edit Member Homestay Fail, Member Homestay not Found",
			ExpectedStatusCode: http.StatusNotFound,
			Uid:                muid,
			Pid:                "",
			In: homestay.EditMemberHomestayIn{
				Name:      "Homestay Name",
				Address:   "Homestay Address",
				Latitude:  "120.12312312",
				Longitude: "90.1212321",
				ImageIds:  []int64{imageId},
			},
		},
		{
			Name:               "Edit Member Homestay Fail, Member Homestay not Found",
			ExpectedStatusCode: http.StatusNotFound,
			Uid:                muid,
			Pid:                "99",
			In: homestay.EditMemberHomestayIn{
				Name:      "Homestay Name",
				Address:   "Homestay Address",
				Latitude:  "120.12312312",
				Longitude: "90.1212321",
				ImageIds:  []int64{imageId},
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
			res := homestayDeps.EditMemberHomestay(ctx, c.Pid, c.Uid, c.In)
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

func TestRemoveMemberHomestay(t *testing.T) {
	err := ClearTables(db)
	if err != nil {
		t.Fatal(err)
	}

	muid, err := createUser(memberRepository, memberSeed)
	if err != nil {
		t.Fatal(err)
	}

	memberHomeId, err := createMemberHomestay(memberHomestayRepository, muid, homestaySeed)
	if err != nil {
		t.Fatal(err)
	}

	mhId := strconv.FormatInt(memberHomeId, 10)

	testCases := []struct {
		ExpectedStatusCode int
		Name               string
		Uid                string
		Pid                string
	}{
		{
			Name:               "Remove Member Homestay Success",
			ExpectedStatusCode: http.StatusOK,
			Uid:                muid,
			Pid:                mhId,
		},
		{
			Name:               "Remove Member Homestay Fail, Member Homestay Not Found",
			ExpectedStatusCode: http.StatusNotFound,
			Uid:                muid,
			Pid:                "",
		},
		{
			Name:               "Remove Member Homestay Fail, Member Homestay Not Found",
			ExpectedStatusCode: http.StatusNotFound,
			Uid:                muid,
			Pid:                "99",
		},
		{
			Name:               "Remove Member Homestay Fail, Member Not Found",
			ExpectedStatusCode: http.StatusNotFound,
			Uid:                "",
			Pid:                mhId,
		},
		{
			Name:               "Remove Member Homestay Fail, Member Not Found",
			ExpectedStatusCode: http.StatusNotFound,
			Uid:                "12345678-1234-1234-1234-123456789012",
			Pid:                mhId,
		},
	}

	for _, c := range testCases {
		t.Run(c.Name, func(t *testing.T) {
			tx, err := db.Begin(context.Background())
			if err != nil {
				t.Fatal(err)
			}

			ctx := context.WithValue(context.Background(), arbitary.TrxX{}, tx)
			res := homestayDeps.RemoveMemberHomestay(ctx, c.Pid, c.Uid)
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

func TestFindMemberHomestay(t *testing.T) {
	err := ClearTables(db)
	if err != nil {
		t.Fatal(err)
	}

	muid, err := createUser(memberRepository, memberSeed)
	if err != nil {
		t.Fatal(err)
	}

	memberHomeId, err := createMemberHomestay(memberHomestayRepository, muid, homestaySeed)
	if err != nil {
		t.Fatal(err)
	}

	mhId := strconv.FormatInt(memberHomeId, 10)

	testCases := []struct {
		ExpectedStatusCode int
		Name               string
		Uid                string
		Pid                string
	}{
		{
			Name:               "Find Member Homestay Success",
			ExpectedStatusCode: http.StatusOK,
			Uid:                muid,
			Pid:                mhId,
		},
		{
			Name:               "Find Member Homestay Fail, Member Homestay Not Found",
			ExpectedStatusCode: http.StatusNotFound,
			Uid:                muid,
			Pid:                "",
		},
		{
			Name:               "Find Member Homestay Fail, Member Homestay Not Found",
			ExpectedStatusCode: http.StatusNotFound,
			Uid:                muid,
			Pid:                "99",
		},
		{
			Name:               "Find Member Homestay Fail, Member Not Found",
			ExpectedStatusCode: http.StatusNotFound,
			Uid:                "",
			Pid:                mhId,
		},
		{
			Name:               "Find Member Homestay Fail, Member Not Found",
			ExpectedStatusCode: http.StatusNotFound,
			Uid:                "12345678-1234-1234-1234-123456789012",
			Pid:                mhId,
		},
	}

	for _, c := range testCases {
		t.Run(c.Name, func(t *testing.T) {
			tx, err := db.Begin(context.Background())
			if err != nil {
				t.Fatal(err)
			}

			ctx := context.WithValue(context.Background(), arbitary.TrxX{}, tx)
			res := homestayDeps.FindMemberHomestay(ctx, c.Pid, c.Uid)
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

func TestQueryMemberHomestay(t *testing.T) {
	err := ClearTables(db)
	if err != nil {
		t.Fatal(err)
	}

	muid, err := createUser(memberRepository, memberSeed)
	if err != nil {
		t.Fatal(err)
	}

	_, err = createMemberHomestay(memberHomestayRepository, muid, homestaySeed)
	if err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		ExpectedStatusCode int
		Name               string
		Uid                string
	}{
		{
			Name:               "Find Member Homestay Success",
			ExpectedStatusCode: http.StatusOK,
			Uid:                muid,
		},
	}

	for _, c := range testCases {
		t.Run(c.Name, func(t *testing.T) {
			tx, err := db.Begin(context.Background())
			if err != nil {
				t.Fatal(err)
			}

			ctx := context.WithValue(context.Background(), arbitary.TrxX{}, tx)
			res := homestayDeps.QueryMemberHomestays(ctx, c.Uid, "", "")
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
