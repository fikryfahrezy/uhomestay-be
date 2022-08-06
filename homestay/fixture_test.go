package homestay_test

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"testing"
	"time"

	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/homestay"
	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/httpdecode"
	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/user"
	"github.com/fikryfahrezy/crypt/agron2"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"golang.org/x/crypto/argon2"
)

var (
	db                       *pgxpool.Pool
	memberRepository         *user.MemberRepository
	homestayImageRepository  *homestay.HomestayImageRepository
	memberHomestayRepository *homestay.MemberHomestayRepository
	homestayDeps             *homestay.HomestayDeps
	fileName                 = "images.jpeg"
	fileDir                  = "./fixture/" + fileName
	fileSeed                 = homestay.HomestayImageModel{
		Name: "file.jpg",
		Url:  "http://localhost:5000/file.jpg",
	}
	homestaySeed = homestay.MemberHomestayModel{
		Name:         "Name",
		Address:      "Address",
		Latitude:     "120",
		Longitude:    "90",
		ThumbnailUrl: "http://localhost:5000/file.jpg",
	}
	memberSeed = user.MemberModel{
		Name:       "Name",
		Username:   "existusername",
		WaPhone:    "+62 821-1111-0000",
		OtherPhone: "+62 821-1111-0000",
		Password:   "password",
		IsAdmin:    true,
		IsApproved: true,
	}
)

var (
	upload homestay.FileUploader = func(filename string, file io.Reader) (string, error) {
		return filename, nil
	}
	captureException homestay.ExceptionCapturer = func(exception error) {}
	captureMessage   homestay.MessageCapturer   = func(message string) {}
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
		`TRUNCATE homestay_images CASCADE`,
		`TRUNCATE member_homestays CASCADE`,
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

func createUser(r *user.MemberRepository, member user.MemberModel) (muid string, err error) {
	memberCp := user.MemberModel(member)

	hash, err := agron2.Argon2Hash(memberCp.Password, "blablabla", 1, 64*1024, 4, 32, argon2.Version, agron2.Argon2Id)
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

func createHomestayImage(r *homestay.HomestayImageRepository, image homestay.HomestayImageModel) (id int64, err error) {
	if image, err = r.Save(context.Background(), image); err != nil {
		return 0, err
	}

	return int64(image.Id), nil
}

func createMemberHomestay(r *homestay.MemberHomestayRepository, memberId string, homestay homestay.MemberHomestayModel) (id int64, err error) {
	homestay.MemberId = memberId
	if homestay, err = r.Save(context.Background(), homestay); err != nil {
		return 0, err
	}

	return int64(homestay.Id), nil
}

func generateFile(fileDir, fileName string) httpdecode.FileHeader {
	f, err := os.OpenFile(fileDir, os.O_RDONLY, 0o444)
	if err != nil {
		log.Fatal(err)
	}

	return httpdecode.FileHeader{
		Filename: fileName,
		File:     f,
	}
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
	homestayImageRepository = homestay.NewHomestayImageRepository(db)
	memberHomestayRepository = homestay.NewMemberHomestayRepository(db)
	homestayDeps = homestay.NewDeps(
		captureMessage,
		captureException,
		upload,
		homestayImageRepository,
		memberHomestayRepository,
		memberRepository,
	)

	if err := LoadTables(db); err != nil {
		log.Fatal(err)
	}

	// Run tests
	code := m.Run()

	// You can't defer this because os.Exit doesn't care for defer
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}
