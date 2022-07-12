package user_test

import (
	"context"
	"embed"
	"fmt"
	"io"
	"log"
	"os"
	"testing"
	"time"

	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/config"
	"github.com/fikryfahrezy/crypt/agron2"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"golang.org/x/crypto/argon2"

	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/user"
)

var (
	db                  *pgxpool.Pool
	memberRepository    *user.MemberRepository
	positionRepository  *user.PositionRepository
	orgRepository       *user.OrgStructureRepository
	orgPeriodRepository *user.OrgPeriodRepository
	goalRepository      *user.GoalRepository
	userDeps            *user.UserDeps
	tmpl                embed.FS
	conf                = config.Config{
		JwtKey:          []byte("testestestest"),
		JwtAudiencesStr: "this",
		JwtKeyStr:       "testestestest",
		JwtIssuerUrl:    "http://localhost:5000",
		Argon2Salt:      "saltingmin8chars",
		JwtAudiences:    []string{"those"},
	}
	fileName = "images.jpeg"
	fileDir  = "./fixture/" + fileName
	period   = user.OrgPeriodModel{
		StartDate: time.Now(),
		EndDate:   time.Now().Add(time.Hour * 24 * 365),
		IsActive:  true,
	}
	inactivePeriod = user.OrgPeriodModel{
		StartDate: time.Now(),
		EndDate:   time.Now().Add(time.Hour * 24 * 365),
		IsActive:  false,
	}
	position = user.PositionModel{
		Name:  "Leader",
		Level: 1,
	}
	member = user.MemberModel{
		Name:              "Name",
		HomestayName:      "Homestay Name",
		Username:          "existusername",
		WaPhone:           "+62 821-1111-9995",
		OtherPhone:        "+62 821-1111-9995",
		HomestayAddress:   "Homestay Address",
		HomestayLatitude:  "120.12312312",
		HomestayLongitude: "90.1212321",
		Password:          "password",
		IsAdmin:           true,
		IsApproved:        true,
	}
	memberAdmin = user.MemberModel{
		Name:              "Name",
		HomestayName:      "Homestay Name",
		Username:          "existusernameone",
		WaPhone:           "+62 821-1111-9996",
		OtherPhone:        "+62 821-1111-9996",
		HomestayAddress:   "Homestay Address",
		HomestayLatitude:  "120.12312312",
		HomestayLongitude: "90.1212321",
		Password:          "password",
		IsAdmin:           true,
		IsApproved:        true,
	}
	memberNormal = user.MemberModel{
		Name:              "Name",
		HomestayName:      "Homestay Name",
		Username:          "existusernametwo",
		WaPhone:           "+62 821-1111-9997",
		OtherPhone:        "+62 821-1111-9997",
		HomestayAddress:   "Homestay Address",
		HomestayLatitude:  "120.12312312",
		HomestayLongitude: "90.1212321",
		Password:          "password",
		IsAdmin:           false,
		IsApproved:        true,
	}
	member2 = user.MemberModel{
		Name:              "Name Two",
		HomestayName:      "Homestay Name Two",
		Username:          "existusernamethree",
		WaPhone:           "+62 821-1111-9998",
		OtherPhone:        "+62 821-1111-9998",
		HomestayAddress:   "Homestay Address Two",
		HomestayLatitude:  "120.12312312",
		HomestayLongitude: "90.1212321",
		Password:          "password",
		IsAdmin:           true,
		IsApproved:        true,
	}
	pendingMember = user.MemberModel{
		Name:              "Name Two",
		HomestayName:      "Homestay Name Two",
		Username:          "existusernamefour",
		WaPhone:           "+62 821-1111-9999",
		OtherPhone:        "+62 821-1111-9999",
		HomestayAddress:   "Homestay Address Two",
		HomestayLatitude:  "120.12312312",
		HomestayLongitude: "90.1212321",
		Password:          "password",
		IsAdmin:           true,
		IsApproved:        false,
	}
	goalSeed = user.GoalModel{
		Vision: map[string]interface{}{
			"test": "test",
		},
		VisionText: "test",
		Mission: map[string]interface{}{
			"test": "test",
		},
		MissionText: "test",
	}
)

var (
	upload user.FileUploader = func(filename string, file io.Reader) (string, error) {
		return "", nil
	}
	captureException user.ExceptionCapturer = func(exception error) {}
	captureMessage   user.MessageCapturer   = func(message string) {}
)

func LoadTables(conn *pgxpool.Pool) error {
	tx, err := conn.Begin(context.Background())
	if err != nil {
		return err
	}

	defer tx.Rollback(context.Background())

	f, err := os.ReadFile("../docs/db.sql")
	if err != nil {
		return err
	}

	_, err = tx.Exec(context.Background(),
		string(f),
	)
	if err != nil {
		return err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func ClearTables(conn *pgxpool.Pool) error {
	tx, err := conn.Begin(context.Background())
	if err != nil {
		return err
	}

	defer tx.Rollback(context.Background())

	// This should be in order of which table truncate first before the other
	queries := []string{
		`TRUNCATE org_structures CASCADE`,
		`TRUNCATE members CASCADE`,
		`TRUNCATE positions CASCADE`,
		`TRUNCATE org_periods CASCADE`,
		`TRUNCATE goals CASCADE`,
	}

	for _, v := range queries {
		_, err = tx.Exec(context.Background(),
			v,
		)
		if err != nil {
			return err
		}
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func createUser(r *user.MemberRepository, member user.MemberModel) (muid string, err error) {
	memberCp := user.MemberModel(member)

	hash, err := agron2.Argon2Hash(memberCp.Password, userDeps.Argon2Salt, 1, 64*1024, 4, 32, argon2.Version, agron2.Argon2Id)
	if err != nil {
		return "", err
	}

	memberCp.Password = hash

	uid, _ := uuid.NewV6()
	memberCp.Id.Scan(uid.String())
	if err := r.Save(context.Background(), memberCp); err != nil {
		return "", err
	}

	return uid.String(), nil
}

func createFullUser(u *user.UserDeps, member user.MemberModel, period user.OrgPeriodModel, position user.PositionModel) (uid string, periodId, positionId uint64, err error) {
	nperiod, err := u.OrgPeriodRepository.Save(context.Background(), period)
	if err != nil {
		return "", 0, 0, err
	}

	nposition, err := u.PositionRepository.Save(context.Background(), position)
	if err != nil {
		return "", 0, 0, err
	}

	uid, err = createUser(u.MemberRepository, member)
	if err != nil {
		return "", 0, 0, err
	}

	orgStructure := user.OrgStructureModel{
		PositionName:  nposition.Name,
		PositionLevel: nposition.Level,
		MemberId:      uid,
		PositionId:    nposition.Id,
		OrgPeriodId:   nperiod.Id,
	}

	if err = u.OrgStructureRepository.Save(context.Background(), orgStructure); err != nil {
		return "", 0, 0, err
	}

	return uid, nperiod.Id, nposition.Id, nil
}

func createOrgStructute(u *user.UserDeps, member user.MemberModel, period user.OrgPeriodModel, position user.PositionModel) (p user.OrgPeriodModel, err error) {
	uid, err := createUser(u.MemberRepository, member)
	if err != nil {
		return user.OrgPeriodModel{}, err
	}

	nperiod, err := u.OrgPeriodRepository.Save(context.Background(), period)
	if err != nil {
		return user.OrgPeriodModel{}, err
	}

	nposition, err := u.PositionRepository.Save(context.Background(), position)
	if err != nil {
		return user.OrgPeriodModel{}, err
	}

	ps := []user.PositionIn{
		{
			Id: int64(nposition.Id),
			Members: []user.MemberIn{
				{
					Id: uid,
				},
			},
		},
	}

	err = u.SaveOrgStructure(context.Background(), nperiod.Id, ps)
	if err != nil {
		return user.OrgPeriodModel{}, err
	}

	return nperiod, nil
}

func TestMain(m *testing.M) {
	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// pulls an image, creates a container based on it and runs it
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "14.1",
		Env: []string{
			"POSTGRES_PASSWORD=secret",
			"POSTGRES_USER=user_name",
			"POSTGRES_DB=dbname",
			"listen_addresses = '*'",
		},
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	hostAndPort := resource.GetHostPort("5432/tcp")
	databaseUrl := fmt.Sprintf("postgres://user_name:secret@%s/dbname?sslmode=disable", hostAndPort)

	log.Println("Connecting to database on url: ", databaseUrl)

	resource.Expire(120) // Tell docker to hard kill the container in 120 seconds

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	pool.MaxWait = 120 * time.Second
	if err = pool.Retry(func() error {
		dbConfig, err := pgxpool.ParseConfig(databaseUrl)
		if err != nil {
			return err
		}

		db, err = pgxpool.ConnectConfig(context.Background(), dbConfig)
		if err != nil {
			return err
		}

		return db.Ping(context.Background())
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	memberRepository = user.NewMemberRepository(db)
	positionRepository = user.NewPositionRepository(db)
	orgRepository = user.NewOrgStructureRepository(db)
	orgPeriodRepository = user.NewOrgPeriodRepository(db)
	goalRepository = user.NewGoalRepository(db)

	userDeps = user.NewDeps(
		conf.JwtKey,
		conf.JwtIssuerUrl,
		conf.Argon2Salt,
		conf.JwtAudiences,
		captureMessage,
		captureException,
		upload,
		tmpl,
		memberRepository,
		positionRepository,
		orgRepository,
		orgPeriodRepository,
		goalRepository,
	)

	LoadTables(db)

	// Run tests
	code := m.Run()

	// You can't defer this because os.Exit doesn't care for defer
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}
