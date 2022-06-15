package dues_test

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"testing"
	"time"

	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/cashflow"
	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/dues"
	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/user"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

var (
	db                   *pgxpool.Pool
	duesRepository       *dues.DuesRepository
	memberDuesRepository *dues.MemberDuesRepository
	memberRepository     *user.MemberRepository
	cashflowRepository   *cashflow.CashflowRepository
	duesDeps             *dues.DuesDeps
	fileName             = "images.jpeg"
	fileDir              = "./fixture/" + fileName
	memberSeed           = user.MemberModel{
		Name:              "Name",
		HomestayName:      "Homestay Name",
		Username:          "existusername",
		WaPhone:           "+62 821-1111-0000",
		OtherPhone:        "+62 821-1111-0000",
		HomestayAddress:   "Homestay Address",
		HomestayLatitude:  "120.12312312",
		HomestayLongitude: "90.1212321",
		Password:          "password",
		IsAdmin:           true,
		IsApproved:        true,
	}
	memberSeed2 = user.MemberModel{
		Name:              "Name Two",
		HomestayName:      "Homestay Name Two",
		Username:          "existusernametwo",
		WaPhone:           "+62 821-1111-0001",
		OtherPhone:        "+62 821-1111-0001",
		HomestayAddress:   "Homestay Address Two",
		HomestayLatitude:  "120.12312312",
		HomestayLongitude: "90.1212321",
		Password:          "password",
		IsAdmin:           true,
		IsApproved:        true,
	}
	duesSeed = dues.DuesModel{
		Date:      time.Now().Add(time.Hour * 750),
		IdrAmount: "20000",
	}
	duesSeed2 = dues.DuesModel{
		Date:      time.Now().Add(time.Hour * 750 * 2),
		IdrAmount: "20000",
	}
	duesSeed3 = dues.DuesModel{
		Date:      time.Now().Add(time.Hour * 750 * 3),
		IdrAmount: "20000",
	}
	pastDuesSeed = dues.DuesModel{
		Date:      time.Now().Add(-1 * (time.Hour * 750)),
		IdrAmount: "20000",
	}
	paidMemDSeed = dues.MemberDuesModel{
		Status: dues.Paid,
	}
	unpaidMemDSeed = dues.MemberDuesModel{
		Status: dues.Unpaid,
	}
)

var upload dues.Uploader = func(filename string, file io.Reader) (string, error) {
	return "", nil
}

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
		`TRUNCATE member_dues CASCADE`,
		`TRUNCATE dues CASCADE`,
		`TRUNCATE members CASCADE`,
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

func createMemberNDues(d *dues.DuesDeps, member user.MemberModel, duesm dues.DuesModel) (muid string, duesId uint64, err error) {
	memberCp := user.MemberModel(member)

	uid, _ := uuid.NewV6()
	memberCp.Id.Scan(uid.String())
	if err := d.MemberRepository.Save(context.Background(), memberCp); err != nil {
		return "", 0, err
	}
	if err != nil {
		return "", 0, err
	}

	nd2, err := d.DuesRepository.Save(context.Background(), duesm)
	if err != nil {
		return "", 0, err
	}

	return uid.String(), nd2.Id, nil
}

func createMemberDues(d *dues.DuesDeps, member user.MemberModel, duesm dues.DuesModel, memberDues dues.MemberDuesModel) (muid string, duesId, memberDuesId uint64, err error) {
	memberCp := user.MemberModel(member)

	uid, _ := uuid.NewV6()
	memberCp.Id.Scan(uid.String())
	if err := d.MemberRepository.Save(context.Background(), memberCp); err != nil {
		return "", 0, 0, err
	}
	if err != nil {
		return "", 0, 0, err
	}

	nd2, err := d.DuesRepository.Save(context.Background(), duesm)
	if err != nil {
		return "", 0, 0, err
	}

	memberDuesCp := dues.MemberDuesModel(memberDues)
	memberDuesCp.MemberId = uid.String()
	memberDuesCp.DuesId = nd2.Id

	nd3, err := d.MemberDuesRepository.Save(context.Background(), memberDuesCp)
	if err != nil {
		return "", 0, 0, err
	}

	return uid.String(), nd2.Id, nd3.Id, nil
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

	duesRepository = dues.NewDeusRepository(db)
	memberDuesRepository = dues.NewMemberDeusRepository(db)
	memberRepository = user.NewMemberRepository(db)
	cashflowRepository = cashflow.NewRepository(db)

	duesDeps = dues.NewDeps(upload, duesRepository, memberDuesRepository, memberRepository, cashflowRepository)

	LoadTables(db)

	// Run tests
	code := m.Run()

	// You can't defer this because os.Exit doesn't care for defer
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}
