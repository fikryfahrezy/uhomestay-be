package user_test

import (
	"context"
	"net/http"
	"strings"
	"testing"

	arbitary "github.com/PA-D3RPLA/d3if43-htt-uhomestay/arbitrary"
	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/user"
	"github.com/stretchr/testify/assert"
	"gopkg.in/guregu/null.v4"
)

func assertUser(t *testing.T, r *user.MemberRepository, u user.AddMemberIn) {
	user, err := r.FindByUsername(u.Username)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, u.Name, user.Name)
	assert.Equal(t, u.OtherPhone, user.OtherPhone)
	assert.Equal(t, u.WaPhone, user.WaPhone)
	assert.Equal(t, u.Username, user.Username)
	assert.Equal(t, u.IsAdmin.Bool, user.IsAdmin)
}

func TestRegisterMember(t *testing.T) {
	err := ClearTables(db)
	if err != nil {
		t.Fatal(err)
	}

	_, err = createUser(memberRepository, member)
	if err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		Name               string
		ExpectedStatusCode int
		In                 user.RegisterIn
	}{
		{
			Name:               "Register Success",
			ExpectedStatusCode: http.StatusCreated,
			In: user.RegisterIn{
				Name:              "Name",
				HomestayName:      "Homestay Name",
				Username:          "username",
				WaPhone:           "+62 821-1111-0001",
				OtherPhone:        "+62 821-1111-0001",
				HomestayAddress:   "Homestay Address",
				HomestayLatitude:  "120.12312312",
				HomestayLongitude: "90.1212321",
				Password:          "password",
				Profile:           generateFile(fileDir, fileName),
				IdCard:            generateFile(fileDir, fileName),
				HomestayPhoto:     generateFile(fileDir, fileName),
			},
		},
		{
			Name:               "Register Fail, Username Exist",
			ExpectedStatusCode: http.StatusBadRequest,
			In: user.RegisterIn{
				Name:              "Name",
				HomestayName:      "Homestay Name",
				Username:          member.Username,
				WaPhone:           "+62 821-1111-0002",
				OtherPhone:        "+62 821-1111-0002",
				HomestayAddress:   "Homestay Address",
				HomestayLatitude:  "120.12312312",
				HomestayLongitude: "90.1212321",
				Password:          "password",
				Profile:           generateFile(fileDir, fileName),
				IdCard:            generateFile(fileDir, fileName),
				HomestayPhoto:     generateFile(fileDir, fileName),
			},
		},
		{
			Name:               "Register Fail, WA Phone Exist",
			ExpectedStatusCode: http.StatusBadRequest,
			In: user.RegisterIn{
				Name:              "Name",
				HomestayName:      "Homestay Name",
				Username:          "username3",
				WaPhone:           member.WaPhone,
				OtherPhone:        "+62 821-1111-0003",
				HomestayAddress:   "Homestay Address",
				HomestayLatitude:  "120.12312312",
				HomestayLongitude: "90.1212321",
				Password:          "password",
				Profile:           generateFile(fileDir, fileName),
				IdCard:            generateFile(fileDir, fileName),
				HomestayPhoto:     generateFile(fileDir, fileName),
			},
		},
		{
			Name:               "Register Fail, Other Phone Exist",
			ExpectedStatusCode: http.StatusBadRequest,
			In: user.RegisterIn{
				Name:              "Name",
				HomestayName:      "Homestay Name",
				Username:          "username3",
				WaPhone:           "+62 821-1111-0004",
				OtherPhone:        member.OtherPhone,
				HomestayAddress:   "Homestay Address",
				HomestayLatitude:  "120.12312312",
				HomestayLongitude: "90.1212321",
				Password:          "password",
				Profile:           generateFile(fileDir, fileName),
				IdCard:            generateFile(fileDir, fileName),
				HomestayPhoto:     generateFile(fileDir, fileName),
			},
		},
		{
			Name:               "Register Fail, Name requied",
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			In: user.RegisterIn{
				Name:              "",
				HomestayName:      "Homestay Name",
				Username:          "username3",
				WaPhone:           "+62 821-1111-0005",
				OtherPhone:        "+62 821-1111-0005",
				HomestayAddress:   "Homestay Address",
				HomestayLatitude:  "120.12312312",
				HomestayLongitude: "90.1212321",
				Password:          "password",
				Profile:           generateFile(fileDir, fileName),
				IdCard:            generateFile(fileDir, fileName),
				HomestayPhoto:     generateFile(fileDir, fileName),
			},
		},
		{
			Name:               "Register Fail, Name over 100 chars",
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			In: user.RegisterIn{
				Name:              strings.Repeat("a", 101),
				HomestayName:      "Homestay Name",
				Username:          "username3",
				WaPhone:           "+62 821-1111-0005",
				OtherPhone:        "+62 821-1111-0005",
				HomestayAddress:   "Homestay Address",
				HomestayLatitude:  "120.12312312",
				HomestayLongitude: "90.1212321",
				Password:          "password",
				Profile:           generateFile(fileDir, fileName),
				IdCard:            generateFile(fileDir, fileName),
				HomestayPhoto:     generateFile(fileDir, fileName),
			},
		},
		{
			Name:               "Register Fail, Wa phone requied",
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			In: user.RegisterIn{
				Name:              "Name",
				HomestayName:      "Homestay Name",
				Username:          "username3",
				WaPhone:           "",
				OtherPhone:        "+62 821-1111-0005",
				HomestayAddress:   "Homestay Address",
				HomestayLatitude:  "120.12312312",
				HomestayLongitude: "90.1212321",
				Password:          "password",
				Profile:           generateFile(fileDir, fileName),
				IdCard:            generateFile(fileDir, fileName),
				HomestayPhoto:     generateFile(fileDir, fileName),
			},
		},
		{
			Name:               "Register Fail, Wa phone over 50 chars",
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			In: user.RegisterIn{
				Name:              "Name",
				HomestayName:      "Homestay Name",
				Username:          "username3",
				WaPhone:           strings.Repeat("0", 51),
				OtherPhone:        "+62 821-1111-0005",
				HomestayAddress:   "Homestay Address",
				HomestayLatitude:  "120.12312312",
				HomestayLongitude: "90.1212321",
				Password:          "password",
				Profile:           generateFile(fileDir, fileName),
				IdCard:            generateFile(fileDir, fileName),
				HomestayPhoto:     generateFile(fileDir, fileName),
			},
		},
		{
			Name:               "Register Fail, other phone requied",
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			In: user.RegisterIn{
				Name:              "",
				HomestayName:      "Homestay Name",
				Username:          "username3",
				WaPhone:           "+62 821-1111-0005",
				OtherPhone:        strings.Repeat("0", 51),
				HomestayAddress:   "Homestay Address",
				HomestayLatitude:  "120.12312312",
				HomestayLongitude: "90.1212321",
				Password:          "password",
				Profile:           generateFile(fileDir, fileName),
				IdCard:            generateFile(fileDir, fileName),
				HomestayPhoto:     generateFile(fileDir, fileName),
			},
		},
		{
			Name:               "Register Fail, Other Phone over 100 chars",
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			In: user.RegisterIn{
				Name:              "Name",
				HomestayName:      "Homestay Name",
				Username:          "username3",
				WaPhone:           "+62 821-1111-0005",
				OtherPhone:        strings.Repeat("0", 51),
				HomestayAddress:   "Homestay Address",
				HomestayLatitude:  "120.12312312",
				HomestayLongitude: "90.1212321",
				Password:          "password",
				Profile:           generateFile(fileDir, fileName),
				IdCard:            generateFile(fileDir, fileName),
				HomestayPhoto:     generateFile(fileDir, fileName),
			},
		},
		{
			Name:               "Register Fail, Homestay Name requied",
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			In: user.RegisterIn{
				Name:              "Name",
				HomestayName:      "",
				Username:          "username3",
				WaPhone:           "+62 821-1111-0005",
				OtherPhone:        "+62 821-1111-0005",
				HomestayAddress:   "Homestay Address",
				HomestayLatitude:  "120.12312312",
				HomestayLongitude: "90.1212321",
				Password:          "password",
				Profile:           generateFile(fileDir, fileName),
				IdCard:            generateFile(fileDir, fileName),
				HomestayPhoto:     generateFile(fileDir, fileName),
			},
		},
		{
			Name:               "Register Fail, Homestay Name over 100 chars",
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			In: user.RegisterIn{
				Name:              "Name",
				HomestayName:      strings.Repeat("0", 101),
				Username:          "username3",
				WaPhone:           "+62 821-1111-0005",
				OtherPhone:        "+62 821-1111-0005",
				HomestayAddress:   "Homestay Address",
				HomestayLatitude:  "120.12312312",
				HomestayLongitude: "90.1212321",
				Password:          "password",
				Profile:           generateFile(fileDir, fileName),
				IdCard:            generateFile(fileDir, fileName),
				HomestayPhoto:     generateFile(fileDir, fileName),
			},
		},
		{
			Name:               "Register Fail, Homestay Address requied",
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			In: user.RegisterIn{
				Name:              "Name",
				HomestayName:      "Homestay Name",
				Username:          "username3",
				WaPhone:           "+62 821-1111-0005",
				OtherPhone:        "+62 821-1111-0005",
				HomestayAddress:   "",
				HomestayLatitude:  "120.12312312",
				HomestayLongitude: "90.1212321",
				Password:          "password",
				Profile:           generateFile(fileDir, fileName),
				IdCard:            generateFile(fileDir, fileName),
				HomestayPhoto:     generateFile(fileDir, fileName),
			},
		},
		{
			Name:               "Register Fail, Homestay Address over 200 chars",
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			In: user.RegisterIn{
				Name:              "Name",
				HomestayName:      "Homestay Name",
				Username:          "username3",
				WaPhone:           "+62 821-1111-0005",
				OtherPhone:        "+62 821-1111-0005",
				HomestayAddress:   strings.Repeat("0", 201),
				HomestayLatitude:  "120.12312312",
				HomestayLongitude: "90.1212321",
				Password:          "password",
				Profile:           generateFile(fileDir, fileName),
				IdCard:            generateFile(fileDir, fileName),
				HomestayPhoto:     generateFile(fileDir, fileName),
			},
		},
		{
			Name:               "Register Fail, Homestay Latitude requied",
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			In: user.RegisterIn{
				Name:              "Name",
				HomestayName:      "Homestay Name",
				Username:          "username3",
				WaPhone:           "+62 821-1111-0005",
				OtherPhone:        "+62 821-1111-0005",
				HomestayAddress:   "Homestay Address",
				HomestayLatitude:  "",
				HomestayLongitude: "90.1212321",
				Password:          "password",
				Profile:           generateFile(fileDir, fileName),
				IdCard:            generateFile(fileDir, fileName),
				HomestayPhoto:     generateFile(fileDir, fileName),
			},
		},
		{
			Name:               "Register Fail, Homestay Latitude over 50 chars",
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			In: user.RegisterIn{
				Name:              "Name",
				HomestayName:      "Homestay Name",
				Username:          "username3",
				WaPhone:           "+62 821-1111-0005",
				OtherPhone:        "+62 821-1111-0005",
				HomestayAddress:   "Homestay Address",
				HomestayLatitude:  strings.Repeat("0", 51),
				HomestayLongitude: "90.1212321",
				Password:          "password",
				Profile:           generateFile(fileDir, fileName),
				IdCard:            generateFile(fileDir, fileName),
				HomestayPhoto:     generateFile(fileDir, fileName),
			},
		},
		{
			Name:               "Register Fail, homestay longitude required",
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			In: user.RegisterIn{
				Name:              "Name",
				HomestayName:      "Homestay Name",
				Username:          "username3",
				WaPhone:           "+62 821-1111-0005",
				OtherPhone:        "+62 821-1111-0005",
				HomestayAddress:   "Homestay Address",
				HomestayLatitude:  "120.12312312",
				HomestayLongitude: "",
				Password:          "password",
				Profile:           generateFile(fileDir, fileName),
				IdCard:            generateFile(fileDir, fileName),
				HomestayPhoto:     generateFile(fileDir, fileName),
			},
		},
		{
			Name:               "Register Fail, homestay longitude over 50 chars",
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			In: user.RegisterIn{
				Name:              "Name",
				HomestayName:      "Homestay Name",
				Username:          "username3",
				WaPhone:           "+62 821-1111-0005",
				OtherPhone:        "+62 821-1111-0005",
				HomestayAddress:   "Homestay Address",
				HomestayLatitude:  "120.12312312",
				HomestayLongitude: strings.Repeat("0", 51),
				Password:          "password",
				Profile:           generateFile(fileDir, fileName),
				IdCard:            generateFile(fileDir, fileName),
				HomestayPhoto:     generateFile(fileDir, fileName),
			},
		},
		{
			Name:               "Register Fail, Username requied",
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			In: user.RegisterIn{
				Name:              "Name",
				HomestayName:      "Homestay Name",
				Username:          "",
				WaPhone:           "+62 821-1111-0005",
				OtherPhone:        "+62 821-1111-0005",
				HomestayAddress:   "Homestay Address",
				HomestayLatitude:  "120.12312312",
				HomestayLongitude: "90.1212321",
				Password:          "password",
				Profile:           generateFile(fileDir, fileName),
				IdCard:            generateFile(fileDir, fileName),
				HomestayPhoto:     generateFile(fileDir, fileName),
			},
		},
		{
			Name:               "Register Fail, Username over 50 chars",
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			In: user.RegisterIn{
				Name:              "Name",
				HomestayName:      "Homestay Name",
				Username:          strings.Repeat("a", 51),
				WaPhone:           "+62 821-1111-0005",
				OtherPhone:        "+62 821-1111-0005",
				HomestayAddress:   "Homestay Address",
				HomestayLatitude:  "120.12312312",
				HomestayLongitude: "90.1212321",
				Password:          "password",
				Profile:           generateFile(fileDir, fileName),
				IdCard:            generateFile(fileDir, fileName),
				HomestayPhoto:     generateFile(fileDir, fileName),
			},
		},
		{
			Name:               "Register Fail, Password required",
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			In: user.RegisterIn{
				Name:              "Name",
				HomestayName:      "Homestay Name",
				Username:          "username3",
				WaPhone:           "+62 821-1111-0005",
				OtherPhone:        "+62 821-1111-0005",
				HomestayAddress:   "Homestay Address",
				HomestayLatitude:  "120.12312312",
				HomestayLongitude: "90.1212321",
				Password:          "",
				Profile:           generateFile(fileDir, fileName),
				IdCard:            generateFile(fileDir, fileName),
				HomestayPhoto:     generateFile(fileDir, fileName),
			},
		},
		{
			Name:               "Register Fail, Password over 200 chars",
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			In: user.RegisterIn{
				Name:              "Name",
				HomestayName:      "Homestay Name",
				Username:          "username3",
				WaPhone:           "+62 821-1111-0005",
				OtherPhone:        "+62 821-1111-0005",
				HomestayAddress:   "Homestay Address",
				HomestayLatitude:  "120.12312312",
				HomestayLongitude: "90.1212321",
				Password:          strings.Repeat("a", 201),
				Profile:           generateFile(fileDir, fileName),
				IdCard:            generateFile(fileDir, fileName),
				HomestayPhoto:     generateFile(fileDir, fileName),
			},
		},
		{
			Name:               "Register Fail, Profile Required",
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			In: user.RegisterIn{
				Name:              "Name",
				HomestayName:      "Homestay Name",
				Username:          "username3",
				WaPhone:           "+62 821-1111-0005",
				OtherPhone:        "+62 821-1111-0005",
				HomestayAddress:   "Homestay Address",
				HomestayLatitude:  "120.12312312",
				HomestayLongitude: "90.1212321",
				Password:          "password",
				IdCard:            generateFile(fileDir, fileName),
				HomestayPhoto:     generateFile(fileDir, fileName),
			},
		},
		{
			Name:               "Register Fail, Profile File Name over 200 characters",
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			In: user.RegisterIn{
				Name:              "Name",
				HomestayName:      "Homestay Name",
				Username:          "username3",
				WaPhone:           "+62 821-1111-0005",
				OtherPhone:        "+62 821-1111-0005",
				HomestayAddress:   "Homestay Address",
				HomestayLatitude:  "120.12312312",
				HomestayLongitude: "90.1212321",
				Password:          "password",
				Profile:           generateFile(fileDir, strings.Repeat("a", 201)),
				IdCard:            generateFile(fileDir, fileName),
				HomestayPhoto:     generateFile(fileDir, fileName),
			},
		},
		{
			Name:               "Register Fail, Profile File not an Image Type",
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			In: user.RegisterIn{
				Name:              "Name",
				HomestayName:      "Homestay Name",
				Username:          "username66",
				WaPhone:           "+62 821-1111-0066",
				OtherPhone:        "+62 821-1111-0066",
				HomestayAddress:   "Homestay Address",
				HomestayLatitude:  "120.12312312",
				HomestayLongitude: "90.1212321",
				Password:          "password",
				Profile:           generateFile("./fixture/pdf.pdf", fileName),
				IdCard:            generateFile(fileDir, fileName),
				HomestayPhoto:     generateFile(fileDir, fileName),
			},
		},
		{
			Name:               "Register Fail, Id Card Required",
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			In: user.RegisterIn{
				Name:              "Name",
				HomestayName:      "Homestay Name",
				Username:          "username3",
				WaPhone:           "+62 821-1111-0005",
				OtherPhone:        "+62 821-1111-0005",
				HomestayAddress:   "Homestay Address",
				HomestayLatitude:  "120.12312312",
				HomestayLongitude: "90.1212321",
				Password:          "password",
				Profile:           generateFile(fileDir, fileName),
				HomestayPhoto:     generateFile(fileDir, fileName),
			},
		},
		{
			Name:               "Register Fail, Id Card File Name over 200 characters",
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			In: user.RegisterIn{
				Name:              "Name",
				HomestayName:      "Homestay Name",
				Username:          "username3",
				WaPhone:           "+62 821-1111-0005",
				OtherPhone:        "+62 821-1111-0005",
				HomestayAddress:   "Homestay Address",
				HomestayLatitude:  "120.12312312",
				HomestayLongitude: "90.1212321",
				Password:          "password",
				Profile:           generateFile(fileDir, fileName),
				IdCard:            generateFile(fileDir, strings.Repeat("a", 201)),
				HomestayPhoto:     generateFile(fileDir, fileName),
			},
		},
		{
			Name:               "Register Fail, Id Card File not an Image Type",
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			In: user.RegisterIn{
				Name:              "Name",
				HomestayName:      "Homestay Name",
				Username:          "username3",
				WaPhone:           "+62 821-1111-0005",
				OtherPhone:        "+62 821-1111-0005",
				HomestayAddress:   "Homestay Address",
				HomestayLatitude:  "120.12312312",
				HomestayLongitude: "90.1212321",
				Password:          "password",
				Profile:           generateFile(fileDir, fileName),
				IdCard:            generateFile("./fixture/pdf.pdf", fileName),
				HomestayPhoto:     generateFile(fileDir, fileName),
			},
		},
		{
			Name:               "Register Fail, Homestay Photo Required",
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			In: user.RegisterIn{
				Name:              "Name",
				HomestayName:      "Homestay Name",
				Username:          "username66",
				WaPhone:           "+62 821-1111-0066",
				OtherPhone:        "+62 821-1111-0066",
				HomestayAddress:   "Homestay Address",
				HomestayLatitude:  "120.12312312",
				HomestayLongitude: "90.1212321",
				Password:          "password",
				Profile:           generateFile(fileDir, fileName),
				IdCard:            generateFile(fileDir, fileName),
			},
		},
		{
			Name:               "Register Fail, Homestay Photo File Name over 200 characters",
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			In: user.RegisterIn{
				Name:              "Name",
				HomestayName:      "Homestay Name",
				Username:          "username3",
				WaPhone:           "+62 821-1111-0005",
				OtherPhone:        "+62 821-1111-0005",
				HomestayAddress:   "Homestay Address",
				HomestayLatitude:  "120.12312312",
				HomestayLongitude: "90.1212321",
				Password:          "password",
				Profile:           generateFile(fileDir, fileName),
				IdCard:            generateFile(fileDir, fileName),
				HomestayPhoto:     generateFile(fileDir, strings.Repeat("a", 201)),
			},
		},
		{
			Name:               "Register Fail, Homestay Photo File not an Image Type",
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			In: user.RegisterIn{
				Name:              "Name",
				HomestayName:      "Homestay Name",
				Username:          "username66",
				WaPhone:           "+62 821-1111-0066",
				OtherPhone:        "+62 821-1111-0066",
				HomestayAddress:   "Homestay Address",
				HomestayLatitude:  "120.12312312",
				HomestayLongitude: "90.1212321",
				Password:          "password",
				Profile:           generateFile(fileDir, fileName),
				IdCard:            generateFile(fileDir, fileName),
				HomestayPhoto:     generateFile("./fixture/pdf.pdf", fileName),
			},
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			tx, err := db.Begin(context.Background())
			if err != nil {
				t.Fatal(err)
			}

			ctx := context.WithValue(context.Background(), arbitary.TrxX{}, tx)
			res := userDeps.MemberRegister(ctx, c.In)
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

func TestAddMember(t *testing.T) {
	err := ClearTables(db)
	if err != nil {
		t.Fatal(err)
	}

	_, err = createUser(memberRepository, member)
	if err != nil {
		t.Fatal(err)
	}

	pr, err := orgPeriodRepository.Save(context.Background(), period)
	if err != nil {
		t.Fatal(err)
	}

	ps, err := positionRepository.Save(context.Background(), position)
	if err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		Name               string
		ExpectedStatusCode int
		In                 user.AddMemberIn
		Assert             func(t *testing.T, r *user.MemberRepository, u user.AddMemberIn)
	}{
		{
			Name:               "Add Member Success",
			ExpectedStatusCode: http.StatusCreated,
			Assert:             assertUser,
			In: user.AddMemberIn{
				Name:        "Name",
				Username:    "username",
				PositionIds: []int64{int64(ps.Id)},
				PeriodId:    int64(pr.Id),
				WaPhone:     "+62 821-1111-0001",
				OtherPhone:  "+62 821-1111-0001",
				Password:    "password",
				IsAdmin:     null.BoolFrom(true),
				Profile:     generateFile(fileDir, fileName),
				IdCard:      generateFile(fileDir, fileName),
			},
		},
		{
			Name:               "Add Member Fail, Username Exist",
			ExpectedStatusCode: http.StatusBadRequest,
			Assert:             func(t *testing.T, r *user.MemberRepository, u user.AddMemberIn) {},
			In: user.AddMemberIn{
				Name:        "Name",
				Username:    member.Username,
				PositionIds: []int64{int64(ps.Id)},
				PeriodId:    int64(pr.Id),
				WaPhone:     "+62 821-1111-0002",
				OtherPhone:  "+62 821-1111-0002",
				Password:    "password",
				IsAdmin:     null.BoolFrom(true),
				Profile:     generateFile(fileDir, fileName),
				IdCard:      generateFile(fileDir, fileName),
			},
		},
		{
			Name:               "Add Member Fail, WA Phone Exist",
			ExpectedStatusCode: http.StatusBadRequest,
			Assert:             func(t *testing.T, r *user.MemberRepository, u user.AddMemberIn) {},
			In: user.AddMemberIn{
				Name:        "Name",
				Username:    "username3",
				PositionIds: []int64{int64(ps.Id)},
				PeriodId:    int64(pr.Id),
				WaPhone:     member.WaPhone,
				OtherPhone:  "+62 821-1111-0003",
				Password:    "password",
				IsAdmin:     null.BoolFrom(true),
				Profile:     generateFile(fileDir, fileName),
				IdCard:      generateFile(fileDir, fileName),
			},
		},
		{
			Name:               "Add Member Fail, Other Phone Exist",
			ExpectedStatusCode: http.StatusBadRequest,
			Assert:             func(t *testing.T, r *user.MemberRepository, u user.AddMemberIn) {},
			In: user.AddMemberIn{
				Name:        "Name",
				Username:    "username3",
				PositionIds: []int64{int64(ps.Id)},
				PeriodId:    int64(pr.Id),
				WaPhone:     "+62 821-1111-0004",
				OtherPhone:  member.OtherPhone,
				Password:    "password",
				IsAdmin:     null.BoolFrom(true),
				Profile:     generateFile(fileDir, fileName),
				IdCard:      generateFile(fileDir, fileName),
			},
		},
		{
			Name:               "Add Member Fail, Position Not Found",
			ExpectedStatusCode: http.StatusNotFound,
			Assert:             func(t *testing.T, r *user.MemberRepository, u user.AddMemberIn) {},
			In: user.AddMemberIn{
				Name:        "Name",
				Username:    "username3",
				PositionIds: []int64{9999},
				PeriodId:    int64(pr.Id),
				WaPhone:     "+62 821-1111-0005",
				OtherPhone:  "+62 821-1111-0005",
				Password:    "password",
				IsAdmin:     null.BoolFrom(true),
				Profile:     generateFile(fileDir, fileName),
				IdCard:      generateFile(fileDir, fileName),
			},
		},
		{
			Name:               "Add Member Fail, Period Not Found",
			ExpectedStatusCode: http.StatusNotFound,
			Assert:             func(t *testing.T, r *user.MemberRepository, u user.AddMemberIn) {},
			In: user.AddMemberIn{
				Name:        "Name",
				Username:    "username3",
				PositionIds: []int64{int64(ps.Id)},
				PeriodId:    9999,
				WaPhone:     "+62 821-1111-0006",
				OtherPhone:  "+62 821-1111-0006",
				Password:    "password",
				IsAdmin:     null.BoolFrom(true),
				Profile:     generateFile(fileDir, fileName),
				IdCard:      generateFile(fileDir, fileName),
			},
		},
		{
			Name:               "Add Member Fail, Name required",
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			Assert:             func(t *testing.T, r *user.MemberRepository, u user.AddMemberIn) {},
			In: user.AddMemberIn{
				Name:        "",
				Username:    "username3",
				WaPhone:     "+62 821-1111-0005",
				OtherPhone:  "+62 821-1111-0005",
				Password:    "password",
				PositionIds: []int64{int64(ps.Id)},
				PeriodId:    9999,
				IsAdmin:     null.BoolFrom(false),
				IdCard:      generateFile(fileDir, fileName),
			},
		},
		{
			Name:               "Add Member Fail, Name over 100 chars",
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			Assert:             func(t *testing.T, r *user.MemberRepository, u user.AddMemberIn) {},
			In: user.AddMemberIn{
				Name:        strings.Repeat("a", 101),
				Username:    "username3",
				WaPhone:     "+62 821-1111-0005",
				OtherPhone:  "+62 821-1111-0005",
				Password:    "password",
				PositionIds: []int64{int64(ps.Id)},
				PeriodId:    9999,
				IsAdmin:     null.BoolFrom(false),
				IdCard:      generateFile(fileDir, fileName),
			},
		},
		{
			Name:               "Add Member Fail, Wa phone requied",
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			Assert:             func(t *testing.T, r *user.MemberRepository, u user.AddMemberIn) {},
			In: user.AddMemberIn{
				Name:        "Name",
				Username:    "username3",
				WaPhone:     "",
				OtherPhone:  "+62 821-1111-0005",
				Password:    "password",
				PositionIds: []int64{int64(ps.Id)},
				PeriodId:    9999,
				IsAdmin:     null.BoolFrom(false),
				IdCard:      generateFile(fileDir, fileName),
			},
		},
		{
			Name:               "Add Member Fail, Wa phone over 50 chars",
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			Assert:             func(t *testing.T, r *user.MemberRepository, u user.AddMemberIn) {},
			In: user.AddMemberIn{
				Name:        "Name",
				Username:    "username3",
				WaPhone:     strings.Repeat("0", 51),
				OtherPhone:  "+62 821-1111-0005",
				Password:    "password",
				PositionIds: []int64{int64(ps.Id)},
				PeriodId:    9999,
				IsAdmin:     null.BoolFrom(false),
				IdCard:      generateFile(fileDir, fileName),
			},
		},
		{
			Name:               "Add Member Fail, Other phone required",
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			Assert:             func(t *testing.T, r *user.MemberRepository, u user.AddMemberIn) {},
			In: user.AddMemberIn{
				Name:        "Name",
				Username:    "username3",
				WaPhone:     "+62 821-1111-0005",
				OtherPhone:  "",
				Password:    "password",
				PositionIds: []int64{int64(ps.Id)},
				PeriodId:    9999,
				IsAdmin:     null.BoolFrom(false),
				IdCard:      generateFile(fileDir, fileName),
			},
		},
		{
			Name:               "Add Member Fail, Other phone over 50 chars",
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			Assert:             func(t *testing.T, r *user.MemberRepository, u user.AddMemberIn) {},
			In: user.AddMemberIn{
				Name:        "Name",
				Username:    "username3",
				WaPhone:     "+62 821-1111-0005",
				OtherPhone:  strings.Repeat("0", 51),
				Password:    "password",
				PositionIds: []int64{int64(ps.Id)},
				PeriodId:    9999,
				IsAdmin:     null.BoolFrom(false),
				IdCard:      generateFile(fileDir, fileName),
			},
		},
		{
			Name:               "Add Member Fail, Username requied",
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			Assert:             func(t *testing.T, r *user.MemberRepository, u user.AddMemberIn) {},
			In: user.AddMemberIn{
				Name:        "Name",
				Username:    "",
				WaPhone:     "+62 821-1111-0005",
				OtherPhone:  "+62 821-1111-0005",
				Password:    "password",
				PositionIds: []int64{int64(ps.Id)},
				PeriodId:    9999,
				IsAdmin:     null.BoolFrom(false),
				IdCard:      generateFile(fileDir, fileName),
			},
		},
		{
			Name:               "Add Member Fail, Username over 50 chars",
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			Assert:             func(t *testing.T, r *user.MemberRepository, u user.AddMemberIn) {},
			In: user.AddMemberIn{
				Name:        "Name",
				Username:    strings.Repeat("a", 51),
				WaPhone:     "+62 821-1111-0005",
				OtherPhone:  "+62 821-1111-0005",
				Password:    "password",
				PositionIds: []int64{int64(ps.Id)},
				PeriodId:    9999,
				IsAdmin:     null.BoolFrom(false),
				IdCard:      generateFile(fileDir, fileName),
			},
		},
		{
			Name:               "Add Member Fail, Password required",
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			Assert:             func(t *testing.T, r *user.MemberRepository, u user.AddMemberIn) {},
			In: user.AddMemberIn{
				Name:        "Name",
				Username:    "username3",
				WaPhone:     "+62 821-1111-0005",
				OtherPhone:  "+62 821-1111-0005",
				Password:    "",
				PositionIds: []int64{int64(ps.Id)},
				PeriodId:    9999,
				IsAdmin:     null.BoolFrom(false),
				IdCard:      generateFile(fileDir, fileName),
			},
		},
		{
			Name:               "Add Member Fail, Password over 200 chars",
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			Assert:             func(t *testing.T, r *user.MemberRepository, u user.AddMemberIn) {},
			In: user.AddMemberIn{
				Name:        "Name",
				Username:    "username3",
				WaPhone:     "+62 821-1111-0005",
				OtherPhone:  "+62 821-1111-0005",
				Password:    strings.Repeat("a", 201),
				PositionIds: []int64{int64(ps.Id)},
				PeriodId:    9999,
				IsAdmin:     null.BoolFrom(false),
				IdCard:      generateFile(fileDir, fileName),
			},
		},
		{
			Name:               "Add Member Fail, Profile not an image",
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			Assert:             func(t *testing.T, r *user.MemberRepository, u user.AddMemberIn) {},
			In: user.AddMemberIn{
				Name:        "Name",
				Username:    "username66",
				PositionIds: []int64{int64(ps.Id)},
				PeriodId:    int64(pr.Id),
				WaPhone:     "+62 821-1111-0066",
				OtherPhone:  "+62 821-1111-0066",
				Password:    "password",
				IsAdmin:     null.BoolFrom(true),
				Profile:     generateFile("./fixture/pdf.pdf", fileName),
			},
		},
		{
			Name:               "Add Member Fail, Id Card not an image",
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			Assert:             func(t *testing.T, r *user.MemberRepository, u user.AddMemberIn) {},
			In: user.AddMemberIn{
				Name:        "Name",
				Username:    "username66",
				PositionIds: []int64{int64(ps.Id)},
				PeriodId:    int64(pr.Id),
				WaPhone:     "+62 821-1111-0066",
				OtherPhone:  "+62 821-1111-0066",
				Password:    "password",
				IsAdmin:     null.BoolFrom(true),
				IdCard:      generateFile("./fixture/pdf.pdf", fileName),
			},
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			tx, err := db.Begin(context.Background())
			if err != nil {
				t.Fatal(err)
			}

			ctx := context.WithValue(context.Background(), arbitary.TrxX{}, tx)
			res := userDeps.AddMember(ctx, c.In)
			tx.Commit(context.Background())
			tx.Rollback(context.Background())

			if res.StatusCode != c.ExpectedStatusCode {
				t.Logf("%#v", res)
				t.Log(c.ExpectedStatusCode)
				t.Fatalf("Expected response code %d. Got %d\n", c.ExpectedStatusCode, res.StatusCode)
			}

			c.Assert(t, memberRepository, c.In)
		})
	}
}

func TestLoginMember(t *testing.T) {
	err := ClearTables(db)
	if err != nil {
		t.Fatal(err)
	}

	_, err = createUser(memberRepository, memberNormal)
	if err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		Name               string
		ExpectedStatusCode int
		In                 user.LoginIn
	}{
		{
			Name:               "Login Member Success",
			ExpectedStatusCode: http.StatusOK,
			In: user.LoginIn{
				Identifier: memberNormal.Username,
				Password:   memberNormal.Password,
			},
		},
		{
			Name:               "Login Member Fail, Wrong Password",
			ExpectedStatusCode: http.StatusBadRequest,
			In: user.LoginIn{
				Identifier: memberNormal.Username,
				Password:   "wrong-password",
			},
		},
		{
			Name:               "Login Member Fail, Username Doesn't Exist",
			ExpectedStatusCode: http.StatusNotFound,
			In: user.LoginIn{
				Identifier: "not-exist",
				Password:   "password",
			},
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			tx, err := db.Begin(context.Background())
			if err != nil {
				t.Fatal(err)
			}

			ctx := context.WithValue(context.Background(), arbitary.TrxX{}, tx)
			res := userDeps.MemberLogin(ctx, c.In)
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

func TestLoginAdmin(t *testing.T) {
	err := ClearTables(db)
	if err != nil {
		t.Fatal(err)
	}

	_, err = createUser(memberRepository, memberAdmin)
	if err != nil {
		t.Fatal(err)
	}

	_, err = createUser(memberRepository, memberNormal)
	if err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		Name               string
		ExpectedStatusCode int
		In                 user.LoginIn
	}{
		{
			Name:               "Login Admin Success",
			ExpectedStatusCode: http.StatusOK,
			In: user.LoginIn{
				Identifier: memberAdmin.Username,
				Password:   memberAdmin.Password,
			},
		},
		{
			Name:               "Login Admin Fail, Wrong Password",
			ExpectedStatusCode: http.StatusBadRequest,
			In: user.LoginIn{
				Identifier: memberAdmin.Username,
				Password:   "wrong-password",
			},
		},
		{
			Name:               "Login Admin Fail, User Not Admin",
			ExpectedStatusCode: http.StatusNotFound,
			In: user.LoginIn{
				Identifier: memberNormal.Username,
				Password:   memberNormal.Password,
			},
		},
		{
			Name:               "Login Admin Fail, Username Doesn't Exist",
			ExpectedStatusCode: http.StatusNotFound,
			In: user.LoginIn{
				Identifier: "not-exist",
				Password:   "password",
			},
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			tx, err := db.Begin(context.Background())
			if err != nil {
				t.Fatal(err)
			}

			ctx := context.WithValue(context.Background(), arbitary.TrxX{}, tx)
			res := userDeps.AdminLogin(ctx, c.In)
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

func TestEditMember(t *testing.T) {
	err := ClearTables(db)
	if err != nil {
		t.Fatal(err)
	}

	uid, prid, psid, err := createFullUser(userDeps, member, period, position)
	if err != nil {
		t.Fatal(err)
	}

	_, _, _, err = createFullUser(userDeps, member2, period, position)
	if err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		Name               string
		ExpectedStatusCode int
		Id                 string
		In                 user.EditMemberIn
	}{
		{
			Name:               "Edit Member Success",
			ExpectedStatusCode: http.StatusOK,
			Id:                 uid,
			In: user.EditMemberIn{
				Name:        "Name",
				Username:    "username",
				PositionIds: []int64{int64(psid)},
				PeriodId:    int64(prid),
				WaPhone:     "+62 821-1111-0003",
				OtherPhone:  "+62 821-1111-0003",
				Password:    "password",
				IsAdmin:     null.BoolFrom(true),
				Profile:     generateFile(fileDir, fileName),
			},
		},
		{
			Name:               "Edit Member Fail, Username Exist",
			ExpectedStatusCode: http.StatusBadRequest,
			Id:                 uid,
			In: user.EditMemberIn{
				Name:        "Name",
				Username:    member2.Username,
				PositionIds: []int64{int64(psid)},
				PeriodId:    int64(prid),
				WaPhone:     "+62 821-1111-0002",
				OtherPhone:  "+62 821-1111-0002",
				Password:    "password",
				IsAdmin:     null.BoolFrom(true),
			},
		},
		{
			Name:               "Edit Member Fail, WA Phone Exist",
			ExpectedStatusCode: http.StatusBadRequest,
			Id:                 uid,
			In: user.EditMemberIn{
				Name:        "Name",
				Username:    "username3",
				PositionIds: []int64{int64(psid)},
				PeriodId:    int64(prid),
				WaPhone:     member2.WaPhone,
				OtherPhone:  "+62 821-1111-0003",
				Password:    "password",
				IsAdmin:     null.BoolFrom(true),
			},
		},
		{
			Name:               "Edit Member Fail, Other Phone Exist",
			ExpectedStatusCode: http.StatusBadRequest,
			Id:                 uid,
			In: user.EditMemberIn{
				Name:        "Name",
				Username:    "username3",
				PositionIds: []int64{int64(psid)},
				PeriodId:    int64(prid),
				WaPhone:     "+62 821-1111-0004",
				OtherPhone:  member2.OtherPhone,
				Password:    "password",
				IsAdmin:     null.BoolFrom(true),
			},
		},
		{
			Name:               "Edit Member Fail, Position Not Found",
			ExpectedStatusCode: http.StatusNotFound,
			Id:                 uid,
			In: user.EditMemberIn{
				Name:        "Name",
				Username:    "username3",
				PositionIds: []int64{9999},
				PeriodId:    int64(prid),
				WaPhone:     "+62 821-1111-0005",
				OtherPhone:  "+62 821-1111-0005",
				Password:    "password",
				IsAdmin:     null.BoolFrom(true),
			},
		},
		{
			Name:               "Edit Member Fail, Period Not Found",
			ExpectedStatusCode: http.StatusNotFound,
			Id:                 uid,
			In: user.EditMemberIn{
				Name:        "Name",
				Username:    "username3",
				PositionIds: []int64{int64(psid)},
				PeriodId:    9999,
				WaPhone:     "+62 821-1111-0006",
				OtherPhone:  "+62 821-1111-0006",
				Password:    "password",
				IsAdmin:     null.BoolFrom(true),
			},
		},
		{
			Name:               "Edit Member Fail, Id not UUID",
			ExpectedStatusCode: http.StatusBadRequest,
			Id:                 "blablabla",
			In: user.EditMemberIn{
				Name:        "Name",
				Username:    "username3",
				PositionIds: []int64{int64(psid)},
				PeriodId:    9999,
				WaPhone:     "+62 821-1111-0006",
				OtherPhone:  "+62 821-1111-0006",
				Password:    "password",
				IsAdmin:     null.BoolFrom(true),
			},
		},
		{
			Name:               "Edit Member Fail, Member Not Found",
			ExpectedStatusCode: http.StatusNotFound,
			Id:                 "12345678-1234-1234-1234-123456789012",
			In: user.EditMemberIn{
				Name:        "Name",
				Username:    "username",
				PositionIds: []int64{int64(psid)},
				PeriodId:    int64(prid),
				WaPhone:     "+62 821-1111-0003",
				OtherPhone:  "+62 821-1111-0003",
				Password:    "password",
				IsAdmin:     null.BoolFrom(true),
			},
		},
		{
			Name:               "Edit Member Fail, Name required",
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			Id:                 "12345678-1234-1234-1234-123456789012",
			In: user.EditMemberIn{
				Name:        "",
				Username:    "username3",
				WaPhone:     "+62 821-1111-0005",
				OtherPhone:  "+62 821-1111-0005",
				Password:    "password",
				PositionIds: []int64{int64(psid)},
				PeriodId:    9999,
				IsAdmin:     null.BoolFrom(true),
			},
		},
		{
			Name:               "Edit Member Fail, Name over 100 chars",
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			Id:                 "12345678-1234-1234-1234-123456789012",
			In: user.EditMemberIn{
				Name:        strings.Repeat("a", 101),
				Username:    "username3",
				WaPhone:     "+62 821-1111-0005",
				OtherPhone:  "+62 821-1111-0005",
				Password:    "password",
				PositionIds: []int64{int64(psid)},
				PeriodId:    9999,
				IsAdmin:     null.BoolFrom(true),
			},
		},
		{
			Name:               "Edit Member Fail, Wa phone required",
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			Id:                 "12345678-1234-1234-1234-123456789012",
			In: user.EditMemberIn{
				Name:        "Name",
				Username:    "username3",
				WaPhone:     "",
				OtherPhone:  "+62 821-1111-0005",
				Password:    "password",
				PositionIds: []int64{int64(psid)},
				PeriodId:    9999,
				IsAdmin:     null.BoolFrom(true),
			},
		},
		{
			Name:               "Edit Member Fail, Wa phone over 50 chars",
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			Id:                 "12345678-1234-1234-1234-123456789012",
			In: user.EditMemberIn{
				Name:        "Name",
				Username:    "username3",
				WaPhone:     strings.Repeat("0", 51),
				OtherPhone:  "+62 821-1111-0005",
				Password:    "password",
				PositionIds: []int64{int64(psid)},
				PeriodId:    9999,
				IsAdmin:     null.BoolFrom(true),
			},
		},
		{
			Name:               "Edit Member Fail, Other phone required",
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			Id:                 "12345678-1234-1234-1234-123456789012",
			In: user.EditMemberIn{
				Name:        "Name",
				Username:    "username3",
				WaPhone:     "+62 821-1111-0005",
				OtherPhone:  "",
				Password:    "password",
				PositionIds: []int64{int64(psid)},
				PeriodId:    9999,
				IsAdmin:     null.BoolFrom(true),
			},
		},
		{
			Name:               "Edit Member Fail, Other phone over 50 chars",
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			Id:                 "12345678-1234-1234-1234-123456789012",
			In: user.EditMemberIn{
				Name:        "Name",
				Username:    "username3",
				WaPhone:     "+62 821-1111-0005",
				OtherPhone:  strings.Repeat("0", 51),
				Password:    "password",
				PositionIds: []int64{int64(psid)},
				PeriodId:    9999,
				IsAdmin:     null.BoolFrom(true),
			},
		},
		{
			Name:               "Edit Member Fail, Username required",
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			Id:                 "12345678-1234-1234-1234-123456789012",
			In: user.EditMemberIn{
				Name:        "Name",
				Username:    "",
				WaPhone:     "+62 821-1111-0005",
				OtherPhone:  "+62 821-1111-0005",
				Password:    "password",
				PositionIds: []int64{int64(psid)},
				PeriodId:    9999,
				IsAdmin:     null.BoolFrom(true),
			},
		},
		{
			Name:               "Edit Member Fail, Username over 50 chars",
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			Id:                 "12345678-1234-1234-1234-123456789012",
			In: user.EditMemberIn{
				Name:        "Name",
				Username:    strings.Repeat("a", 51),
				WaPhone:     "+62 821-1111-0005",
				OtherPhone:  "+62 821-1111-0005",
				Password:    "password",
				PositionIds: []int64{int64(psid)},
				PeriodId:    9999,
				IsAdmin:     null.BoolFrom(true),
			},
		},
		{
			Name:               "Edit Member Fail, Password over 200 chars",
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			Id:                 "12345678-1234-1234-1234-123456789012",
			In: user.EditMemberIn{
				Name:        "Name",
				Username:    "username3",
				WaPhone:     "+62 821-1111-0005",
				OtherPhone:  "+62 821-1111-0005",
				Password:    strings.Repeat("a", 201),
				PositionIds: []int64{int64(psid)},
				PeriodId:    9999,
				IsAdmin:     null.BoolFrom(true),
			},
		},
		{
			Name:               "Edit Member Fail, Profile filename over 200 characters",
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			Id:                 uid,
			In: user.EditMemberIn{
				Name:        "Name",
				Username:    "username66",
				PositionIds: []int64{int64(psid)},
				PeriodId:    int64(prid),
				WaPhone:     "+62 821-1111-0066",
				OtherPhone:  "+62 821-1111-0066",
				Password:    "password",
				IsAdmin:     null.BoolFrom(true),
				Profile:     generateFile(fileDir, strings.Repeat("a", 201)),
			},
		},
		{
			Name:               "Edit Member Fail, Profile not an image",
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			Id:                 uid,
			In: user.EditMemberIn{
				Name:        "Name",
				Username:    "username66",
				PositionIds: []int64{int64(psid)},
				PeriodId:    int64(prid),
				WaPhone:     "+62 821-1111-0066",
				OtherPhone:  "+62 821-1111-0066",
				Password:    "password",
				IsAdmin:     null.BoolFrom(true),
				Profile:     generateFile("./fixture/pdf.pdf", fileName),
			},
		},
		{
			Name:               "Edit Member Fail, Id Card filename over 200 characters",
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			Id:                 uid,
			In: user.EditMemberIn{
				Name:        "Name",
				Username:    "username66",
				PositionIds: []int64{int64(psid)},
				PeriodId:    int64(prid),
				WaPhone:     "+62 821-1111-0066",
				OtherPhone:  "+62 821-1111-0066",
				Password:    "password",
				IsAdmin:     null.BoolFrom(true),
				IdCard:      generateFile(fileDir, strings.Repeat("a", 201)),
			},
		},
		{
			Name:               "Edit Member Fail, Id Card not an image",
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			Id:                 uid,
			In: user.EditMemberIn{
				Name:        "Name",
				Username:    "username66",
				PositionIds: []int64{int64(psid)},
				PeriodId:    int64(prid),
				WaPhone:     "+62 821-1111-0066",
				OtherPhone:  "+62 821-1111-0066",
				Password:    "password",
				IsAdmin:     null.BoolFrom(true),
				IdCard:      generateFile("./fixture/pdf.pdf", fileName),
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
			res := userDeps.EditMember(ctx, c.Id, c.In)
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

func TestRemoveMember(t *testing.T) {
	err := ClearTables(db)
	if err != nil {
		t.Fatal(err)
	}

	uid, err := createUser(memberRepository, member)
	if err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		Name               string
		ExpectedStatusCode int
		Id                 string
	}{
		{
			Name:               "Remove Member Success",
			ExpectedStatusCode: http.StatusOK,
			Id:                 uid,
		},
		{
			Name:               "Remove Member Fail, ID not UUID",
			ExpectedStatusCode: http.StatusNotFound,
			Id:                 "blablabla",
		},
		{
			Name:               "Remove Member Fail, Member Not Found",
			ExpectedStatusCode: http.StatusNotFound,
			Id:                 "12345678-1234-1234-1234-123456789012",
		},
	}

	for _, c := range testCases {
		t.Run(c.Name, func(t *testing.T) {
			tx, err := db.Begin(context.Background())
			if err != nil {
				t.Fatal(err)
			}

			ctx := context.WithValue(context.Background(), arbitary.TrxX{}, tx)
			res := userDeps.RemoveMember(ctx, c.Id)
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

func TestQueryMember(t *testing.T) {
	err := ClearTables(db)
	if err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		Name               string
		ExpectedStatusCode int
	}{
		{
			Name:               "Query Member Success",
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
			res := userDeps.QueryMember(ctx, "", "", "0")
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

func TestFindMemberDetail(t *testing.T) {
	err := ClearTables(db)
	if err != nil {
		t.Fatal(err)
	}

	uid, err := createUser(memberRepository, member)
	if err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		Name               string
		ExpectedStatusCode int
		Id                 string
	}{
		{
			Name:               "Find Member Detail Success",
			ExpectedStatusCode: http.StatusOK,
			Id:                 uid,
		},
		{
			Name:               "Find Member Detail Fail, ID not UUID",
			ExpectedStatusCode: http.StatusNotFound,
			Id:                 "blablabla",
		},
		{
			Name:               "Find Member Detail Fail, Member Not Found",
			ExpectedStatusCode: http.StatusNotFound,
			Id:                 "12345678-1234-1234-1234-123456789012",
		},
	}

	for _, c := range testCases {
		t.Run(c.Name, func(t *testing.T) {
			res := userDeps.FindMemberDetail(context.Background(), c.Id)

			if res.StatusCode != c.ExpectedStatusCode {
				t.Logf("%#v", res)
				t.Log(err)
				t.Fatalf("Expected response code %d. Got %d\n", c.ExpectedStatusCode, res.StatusCode)
			}
		})
	}
}

func TestApproveMember(t *testing.T) {
	err := ClearTables(db)
	if err != nil {
		t.Fatal(err)
	}

	uid1, err := createUser(memberRepository, member)
	if err != nil {
		t.Fatal(err)
	}

	uid2, err := createUser(memberRepository, pendingMember)
	if err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		Name               string
		ExpectedStatusCode int
		Id                 string
	}{
		{
			Name:               "Approve Member Success",
			ExpectedStatusCode: http.StatusOK,
			Id:                 uid2,
		},
		{
			Name:               "Approve Member Fail, ID not UUID",
			ExpectedStatusCode: http.StatusNotFound,
			Id:                 "blablabla",
		},
		{
			Name:               "Approve Member Fail, Member Not Found",
			ExpectedStatusCode: http.StatusNotFound,
			Id:                 "12345678-1234-1234-1234-123456789012",
		},
		{
			Name:               "Approve Member Fail, Member Already Approved",
			ExpectedStatusCode: http.StatusNotFound,
			Id:                 uid1,
		},
	}

	for _, c := range testCases {
		t.Run(c.Name, func(t *testing.T) {
			tx, err := db.Begin(context.Background())
			if err != nil {
				t.Fatal(err)
			}

			ctx := context.WithValue(context.Background(), arbitary.TrxX{}, tx)
			res := userDeps.ApproveMember(ctx, c.Id)
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

func TestUpdateProfile(t *testing.T) {
	err := ClearTables(db)
	if err != nil {
		t.Fatal(err)
	}

	uid, _, _, err := createFullUser(userDeps, member, period, position)
	if err != nil {
		t.Fatal(err)
	}

	_, _, _, err = createFullUser(userDeps, member2, period, position)
	if err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		Name               string
		ExpectedStatusCode int
		Id                 string
		In                 user.UpdateProfileIn
	}{
		{
			Name:               "Update Profile Success",
			ExpectedStatusCode: http.StatusOK,
			Id:                 uid,
			In: user.UpdateProfileIn{
				Name:       "Name",
				Username:   "username",
				WaPhone:    "+62 821-1111-0003",
				OtherPhone: "+62 821-1111-0003",
				Password:   "password",
			},
		},
		{
			Name:               "Update Profile Fail, Username Exist",
			ExpectedStatusCode: http.StatusBadRequest,
			Id:                 uid,
			In: user.UpdateProfileIn{
				Name:       "Name",
				Username:   member2.Username,
				WaPhone:    "+62 821-1111-0002",
				OtherPhone: "+62 821-1111-0002",
				Password:   "password",
			},
		},
		{
			Name:               "Update Profile Fail, WA Phone Exist",
			ExpectedStatusCode: http.StatusBadRequest,
			Id:                 uid,
			In: user.UpdateProfileIn{
				Name:       "Name",
				Username:   "username3",
				WaPhone:    member2.WaPhone,
				OtherPhone: "+62 821-1111-0003",
				Password:   "password",
			},
		},
		{
			Name:               "Update Profile Fail, Other Phone Exist",
			ExpectedStatusCode: http.StatusBadRequest,
			Id:                 uid,
			In: user.UpdateProfileIn{
				Name:       "Name",
				Username:   "username3",
				WaPhone:    "+62 821-1111-0004",
				OtherPhone: member2.OtherPhone,
				Password:   "password",
			},
		},
		{
			Name:               "Update Profile Fail, Id not UUID",
			ExpectedStatusCode: http.StatusNotFound,
			Id:                 "blablabla",
			In: user.UpdateProfileIn{
				Name:       "Name",
				Username:   "username3",
				WaPhone:    "+62 821-1111-0006",
				OtherPhone: "+62 821-1111-0006",
				Password:   "password",
			},
		},
		{
			Name:               "Update Profile Fail, Member Not Found",
			ExpectedStatusCode: http.StatusNotFound,
			Id:                 "12345678-1234-1234-1234-123456789012",
			In: user.UpdateProfileIn{
				Name:       "Name",
				Username:   "username",
				WaPhone:    "+62 821-1111-0003",
				OtherPhone: "+62 821-1111-0003",
				Password:   "password",
			},
		},
		{
			Name:               "Update Profile Fail, Name over 100 chars",
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			Id:                 "12345678-1234-1234-1234-123456789012",
			In: user.UpdateProfileIn{
				Name:       strings.Repeat("a", 101),
				Username:   "username3",
				WaPhone:    "+62 821-1111-0005",
				OtherPhone: "+62 821-1111-0005",
				Password:   "password",
			},
		},
		{
			Name:               "Update Profile Fail, Wa phone over 50 chars",
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			Id:                 "12345678-1234-1234-1234-123456789012",
			In: user.UpdateProfileIn{
				Name:       "Name",
				Username:   "username3",
				WaPhone:    strings.Repeat("0", 51),
				OtherPhone: "+62 821-1111-0005",
				Password:   "password",
			},
		},
		{
			Name:               "Update Profile Fail, Name over 100 chars",
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			Id:                 "12345678-1234-1234-1234-123456789012",
			In: user.UpdateProfileIn{
				Name:       "Name",
				Username:   "username3",
				WaPhone:    "+62 821-1111-0005",
				OtherPhone: strings.Repeat("0", 51),
				Password:   "password",
			},
		},
		{
			Name:               "Update Profile Fail, Username over 50 chars",
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			Id:                 "12345678-1234-1234-1234-123456789012",
			In: user.UpdateProfileIn{
				Name:       "Name",
				Username:   strings.Repeat("a", 51),
				WaPhone:    "+62 821-1111-0005",
				OtherPhone: "+62 821-1111-0005",
				Password:   "password",
			},
		},
		{
			Name:               "Update Profile Fail, Password over 200 chars",
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			Id:                 "12345678-1234-1234-1234-123456789012",
			In: user.UpdateProfileIn{
				Name:       "Name",
				Username:   "username3",
				WaPhone:    "+62 821-1111-0005",
				OtherPhone: "+62 821-1111-0005",
				Password:   strings.Repeat("a", 201),
			},
		},
		{
			Name:               "Update Profile Fail, Profile filename over 200 characters",
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			Id:                 uid,
			In: user.UpdateProfileIn{
				Name:       "Name",
				Username:   "username66",
				WaPhone:    "+62 821-1111-0066",
				OtherPhone: "+62 821-1111-0066",
				Password:   "password",
				Profile:    generateFile(fileDir, strings.Repeat("a", 201)),
			},
		},
		{
			Name:               "Update Profile Fail, Profile not an image",
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			Id:                 uid,
			In: user.UpdateProfileIn{
				Name:       "Name",
				Username:   "username66",
				WaPhone:    "+62 821-1111-0066",
				OtherPhone: "+62 821-1111-0066",
				Password:   "password",
				Profile:    generateFile("./fixture/pdf.pdf", fileName),
			},
		},
		{
			Name:               "Update Profile Fail, Id Card filename over 200 characters",
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			Id:                 uid,
			In: user.UpdateProfileIn{
				Name:       "Name",
				Username:   "username66",
				WaPhone:    "+62 821-1111-0066",
				OtherPhone: "+62 821-1111-0066",
				Password:   "password",
				IdCard:     generateFile(fileDir, strings.Repeat("a", 201)),
			},
		},
		{
			Name:               "Update Profile Fail, Id Card not an image",
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			Id:                 uid,
			In: user.UpdateProfileIn{
				Name:       "Name",
				Username:   "username66",
				WaPhone:    "+62 821-1111-0066",
				OtherPhone: "+62 821-1111-0066",
				Password:   "password",
				IdCard:     generateFile("./fixture/pdf.pdf", fileName),
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
			res := userDeps.UpdateProfile(ctx, c.Id, c.In)
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
