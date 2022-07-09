package blog_test

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"testing"
	"time"

	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/blog"
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

var (
	postgrePool    *pgxpool.Pool
	redisClient    *redis.Client
	blogRepository *blog.BlogRepository
	blogDeps       *blog.BlogDeps
	fileName       = "images.jpeg"
	fileDir        = "./fixture/" + fileName
	imgTmpFolder   = "blabla"
	imgFolder      = "blublu"
	blogSeed       = blog.BlogModel{
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
	upload blog.FileUploader = func(filename string, file io.Reader) (string, string, error) {
		return "", "", nil
	}
	moveFile blog.FileMover = func(from, to string) (string, error) {
		return "", nil
	}
	captureException blog.ExceptionCapturer = func(exception error) {
	}
	captureMessage blog.MessageCapturer = func(message string) {
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
		`TRUNCATE blogs CASCADE`,
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

func ClearRedis(client *redis.Client) error {
	_, err := client.FlushDB(context.Background()).Result()
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

	redisResource, err := pool.Run("redis", "7.0.0", nil)
	if err != nil {
		log.Fatalf("Could not start redis resource: %s", err)
	}

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	err = pool.Retry(func() error {
		var err error
		redisClient = redis.NewClient(&redis.Options{
			Addr: fmt.Sprintf("localhost:%s", redisResource.GetPort("6379/tcp")),
		})

		err = redisClient.Ping(context.TODO()).Err()
		if err != nil {
			return err
		}

		dbConfig, err := pgxpool.ParseConfig(databaseUrl)
		if err != nil {
			return err
		}

		postgrePool, err = pgxpool.ConnectConfig(context.Background(), dbConfig)
		if err != nil {
			return err
		}

		return postgrePool.Ping(context.Background())
	})

	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	blogRepository = blog.NewRepository("imgchc", redisClient, postgrePool)
	blogDeps = blog.NewDeps(
		imgFolder,
		imgTmpFolder,
		captureMessage,
		captureException,
		moveFile,
		upload,
		blogRepository,
	)

	LoadTables(postgrePool)

	// run tests
	code := m.Run()

	// When you're done, kill and remove the container
	if err = pool.Purge(postgreResource); err != nil {
		log.Fatalf("Could not purge mongo resource: %s", err)
	}

	if err = pool.Purge(redisResource); err != nil {
		log.Fatalf("Could not purge redis resource: %s", err)
	}

	if err = redisClient.Close(); err != nil {
		panic(err)
	}

	os.Exit(code)
}
