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
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	mongoClient    *mongo.Client
	redisClient    *redis.Client
	blogRepository *blog.BlogRepository
	blogDeps       *blog.BlogDeps

	fileName     = "images.jpeg"
	fileDir      = "./fixture/" + fileName
	mongoDbName  = "uhomestay"
	imgTmpFolder = "blabla"
	imgFolder    = "blublu"

	blogSeed = blog.BlogModel{
		Title:        "title",
		ShortDesc:    "Short desc",
		Slug:         "slug",
		ThumbnailUrl: "http://localhost:8080/images.jpg",
		Content: map[string]interface{}{
			"test": "hi",
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
)

var (
	upload blog.Uploader = func(filename string, file io.Reader) (string, string, error) {
		return "", "", nil
	}
	moveFile blog.Mover = func(from, to string) (string, error) {
		return "", nil
	}
)

func ClearMongo(client *mongo.Client) error {
	err := client.Database(mongoDbName).Drop(context.Background())
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

	// pull mongodb docker image for version 5.0
	mongoResource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "mongo",
		Tag:        "5.0.6",
		Env: []string{
			// username and password for mongodb superuser
			"MONGO_INITDB_ROOT_USERNAME=root",
			"MONGO_INITDB_ROOT_PASSWORD=password",
		},
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
	})
	if err != nil {
		log.Fatalf("Could not start mongo resource: %s", err)
	}

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

		mongoClient, err = mongo.Connect(
			context.TODO(),
			options.Client().ApplyURI(
				fmt.Sprintf("mongodb://root:password@localhost:%s", mongoResource.GetPort("27017/tcp")),
			),
		)

		err = mongoClient.Ping(context.TODO(), nil)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	blogRepository = blog.NewRepository(mongoDbName, "blogs", "imgchc", redisClient, mongoClient)
	blogDeps = blog.NewDeps(imgFolder, imgTmpFolder, moveFile, upload, blogRepository)

	// run tests
	code := m.Run()

	// When you're done, kill and remove the container
	if err = pool.Purge(mongoResource); err != nil {
		log.Fatalf("Could not purge mongo resource: %s", err)
	}

	if err = pool.Purge(redisResource); err != nil {
		log.Fatalf("Could not purge redis resource: %s", err)
	}

	// disconnect mongodb client
	if err = mongoClient.Disconnect(context.TODO()); err != nil {
		panic(err)
	}

	if err = redisClient.Close(); err != nil {
		panic(err)
	}

	os.Exit(code)
}
