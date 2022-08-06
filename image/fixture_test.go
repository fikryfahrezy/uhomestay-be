package image_test

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"testing"
	"time"

	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/httpdecode"
	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/image"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

var (
	db              *pgxpool.Pool
	imageRepository *image.ImageRepository
	imageDeps       *image.ImageDeps
	fileName        = "images.jpeg"
	fileDir         = "./fixture/" + fileName
	fileSeed        = image.ImageModel{
		Name: "file.jpg",
		Url:  "http://localhost:5000/file.jpg",
	}
)

var (
	upload image.FileUploader = func(filename string, file io.Reader) (string, error) {
		return filename, nil
	}
	captureException image.ExceptionCapturer = func(exception error) {}
	captureMessage   image.MessageCapturer   = func(message string) {}
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
		`TRUNCATE images CASCADE`,
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

	imageRepository = image.NewRepository(db)
	imageDeps = image.NewDeps(
		captureMessage,
		captureException,
		upload,
		imageRepository,
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
