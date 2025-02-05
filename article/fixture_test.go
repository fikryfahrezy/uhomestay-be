package article_test

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"testing"
	"time"

	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/article"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

var (
	postgrePool       *pgxpool.Pool
	articleRepository *article.ArticleRepository
	articleDeps       *article.ArticleDeps
	fileName          = "images.jpeg"
	fileDir           = "./fixture/" + fileName
	imgTmpFolder      = "blabla"
	imgFolder         = "blublu"
	articleSeed       = article.ArticleModel{
		Title:        "title",
		ShortDesc:    "Short desc",
		Slug:         "slug",
		ThumbnailUrl: "http://localhost:8080/images.jpg",
		Content: map[string]interface{}{
			"test": "hi",
		},
		ContentText: "hi",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
)

var (
	upload article.FileUploader = func(filename string, file io.Reader) (string, string, error) {
		return filename, filename, nil
	}
	moveFile article.FileMover = func(from, to string) (string, error) {
		return "", nil
	}
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
		`TRUNCATE articles CASCADE`,
		`TRUNCATE image_caches CASCADE`,
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
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// pulls an image, creates a container based on it and runs it
	postgreResource, err := pool.RunWithOptions(&dockertest.RunOptions{
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
		log.Fatalf("Could not start postgre resource: %s", err)
	}

	hostAndPort := postgreResource.GetHostPort("5432/tcp")
	databaseUrl := fmt.Sprintf("postgres://user_name:secret@%s/dbname?sslmode=disable", hostAndPort)

	log.Println("Connecting to postgre database on url: ", databaseUrl)

	postgreResource.Expire(120) // Tell docker to hard kill the container in 120 seconds

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	pool.MaxWait = 120 * time.Second
	if err = pool.Retry(func() error {
		dbConfig, err := pgxpool.ParseConfig(databaseUrl)
		if err != nil {
			return err
		}

		postgrePool, err = pgxpool.ConnectConfig(context.Background(), dbConfig)
		if err != nil {
			return err
		}

		return postgrePool.Ping(context.Background())
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	articleRepository = article.NewRepository("imgchc", postgrePool)
	articleDeps = article.NewDeps(
		imgFolder,
		imgTmpFolder,
		moveFile,
		upload,
		articleRepository,
	)

	LoadTables(postgrePool)

	// run tests
	code := m.Run()

	// When you're done, kill and remove the container
	if err = pool.Purge(postgreResource); err != nil {
		log.Fatalf("Could not purge postgre resource: %s", err)
	}

	os.Exit(code)
}
