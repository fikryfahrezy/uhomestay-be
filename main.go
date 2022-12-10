// package main

// import (
// 	"context"
// 	"embed"
// 	"log"
// 	"time"

// 	"go.mongodb.org/mongo-driver/mongo"
// 	"go.mongodb.org/mongo-driver/mongo/options"
// 	"go.mongodb.org/mongo-driver/mongo/readpref"

// 	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/article"
// 	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/cashflow"
// 	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/config"
// 	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/dashboard"
// 	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/document"
// 	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/dues"
// 	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/handler"
// 	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/history"
// 	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/homestay"
// 	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/image"
// 	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/user"

// 	"github.com/cloudinary/cloudinary-go"
// 	"github.com/cloudinary/cloudinary-go/api/uploader"
// 	"github.com/getsentry/sentry-go"
// 	"github.com/go-redis/redis/v8"
// 	"github.com/jackc/pgtype"
// 	pgtypeuuid "github.com/jackc/pgtype/ext/gofrs-uuid"
// 	"github.com/jackc/pgx/v4"
// 	"github.com/jackc/pgx/v4/pgxpool"
// 	"github.com/logdna/logdna-go/logger"
// )

// //go:embed tmpl/*
// var tmpl embed.FS

// var (
// 	buildDate  string = "N/A"
// 	commitHash string = "N/A"
// )

// func main() {
// 	conf := config.LoadConfig()

// 	posgreConfig, err := pgxpool.ParseConfig(conf.PostgreUrl)
// 	if err != nil {
// 		log.Fatalf("fail parse posgre database config: %s", err)
// 	}
// 	posgreConfig.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
// 		conn.ConnInfo().RegisterDataType(pgtype.DataType{
// 			Value: &pgtypeuuid.UUID{},
// 			Name:  "uuid",
// 			OID:   pgtype.UUIDOID,
// 		})
// 		return nil
// 	}

// 	posgrePool, err := pgxpool.ConnectConfig(context.Background(), posgreConfig)
// 	if err != nil {
// 		log.Fatalf("fail connect to postgre: %s", err)
// 	}
// 	defer posgrePool.Close()

// 	if _, err := posgrePool.Exec(context.Background(), "SELECT 1;"); err != nil {
// 		log.Fatalf("fail to ping postgre: %s", err)
// 		return
// 	}

// 	// Create a new client and connect to the server
// 	mongoClient, err := mongo.Connect(context.Background(), options.Client().ApplyURI(conf.MongoUri))
// 	if err != nil {
// 		log.Fatalf("fail to connect mongodb: %s", err)
// 	}
// 	defer func() {
// 		if err = mongoClient.Disconnect(context.Background()); err != nil {
// 			log.Fatalf("fail to disconnect mongodb: %s", err)
// 		}
// 	}()
// 	// Ping the primary
// 	if err := mongoClient.Ping(context.Background(), readpref.Primary()); err != nil {
// 		log.Fatalf("fail to ping mongodb: %s", err)
// 	}

// 	opt, err := redis.ParseURL(conf.RedisUrl)
// 	if err != nil {
// 		log.Fatalf("fail to parse redis url: %s", err)
// 	}
// 	redisClient := redis.NewClient(opt)
// 	defer redisClient.Close()

// 	if err := redisClient.Ping(context.Background()).Err(); err != nil {
// 		log.Fatalf("fail to ping redis: %s", err)
// 	}

// 	cld, err := cloudinary.NewFromURL(conf.CloudinaryUrl)
// 	if err != nil {
// 		log.Fatalf("cloudinary.NewFromURL: %s", err)
// 	}

// 	_, err = logger.NewLogger(logger.Options{App: "U-Homestay"}, conf.LogDnaKey)
// 	if err != nil {
// 		log.Fatalf("LogDNA.NewLogger: %s", err)
// 	}

// 	sentrySampleRate := 0.7
// 	if conf.Env == "dev" {
// 		sentrySampleRate = 0
// 	}

// 	err = sentry.Init(sentry.ClientOptions{
// 		Dsn:              conf.SentryDsn,
// 		Environment:      "all",
// 		Release:          "uhomestay@1.0.0",
// 		Debug:            true,
// 		TracesSampleRate: sentrySampleRate,
// 	})
// 	if err != nil {
// 		log.Fatalf("sentry.Init: %s", err)
// 	}

// 	// Flush buffered events before the program terminates.
// 	defer sentry.Flush(2 * time.Second)

// 	memberRepository := user.NewMemberRepository(posgrePool)
// 	positionRepository := user.NewPositionRepository(posgrePool)
// 	orgRepository := user.NewOrgStructureRepository(posgrePool)
// 	periodRepository := user.NewOrgPeriodRepository(posgrePool)
// 	goalRepository := user.NewGoalRepository(posgrePool)
// 	documentRepository := document.NewRepository(posgrePool)
// 	cashflowRepository := cashflow.NewRepository(posgrePool)
// 	duesRepository := dues.NewDeusRepository(posgrePool)
// 	memberDuesRepository := dues.NewMemberDeusRepository(posgrePool)
// 	imageRepository := image.NewRepository(posgrePool)

// 	historyRepository := history.NewRepository(
// 		posgrePool,
// 	)
// 	articleRepository := article.NewRepository(
// 		"imgchc",
// 		redisClient,
// 		posgrePool,
// 	)

// 	memberHomestayRepository := homestay.NewMemberHomestayRepository(
// 		posgrePool,
// 	)
// 	homestayImageRepository := homestay.NewHomestayImageRepository(
// 		posgrePool,
// 	)

// 	userDeps := user.NewDeps(
// 		conf.JwtKey,
// 		conf.JwtIssuerUrl,
// 		conf.Argon2Salt,
// 		conf.JwtAudiences,
// 		user.CaptureMessage(sentry.CaptureMessage),
// 		user.CaptureExeption(sentry.CaptureException),
// 		user.FileUpload(uploader.UploadParams{
// 			Transformation: "c_crop,g_center/q_auto/f_auto",
// 			Tags:           []string{"profile"},
// 			Folder:         "uhomestay/profile",
// 			ResourceType:   "image",
// 		}, cld.Upload.Upload),
// 		tmpl,
// 		memberRepository,
// 		positionRepository,
// 		orgRepository,
// 		periodRepository,
// 		goalRepository,
// 	)

// 	documentDeps := document.NewDeps(
// 		document.CaptureMessage(sentry.CaptureMessage),
// 		document.CaptureExeption(sentry.CaptureException),
// 		document.FileUpload(uploader.UploadParams{
// 			Tags:         []string{"document"},
// 			Folder:       "uhomestay/document",
// 			ResourceType: "raw",
// 		}, cld.Upload.Upload),
// 		documentRepository,
// 	)

// 	historyDeps := history.NewDeps(
// 		history.CaptureMessage(sentry.CaptureMessage),
// 		history.CaptureExeption(sentry.CaptureException),
// 		historyRepository,
// 	)

// 	articleImgFolder := "uhomestay/blog-images-tmp"
// 	articleDeps := article.NewDeps(
// 		"uhomestay/blog-images",
// 		articleImgFolder,
// 		article.CaptureMessage(sentry.CaptureMessage),
// 		article.CaptureExeption(sentry.CaptureException),
// 		article.FileMove(cld.Upload.Rename),
// 		article.FileUpload(uploader.UploadParams{
// 			Tags:         []string{"blogs"},
// 			Folder:       articleImgFolder,
// 			ResourceType: "raw",
// 		}, cld.Upload.Upload),
// 		articleRepository,
// 	)

// 	cashflowDeps := cashflow.NewDeps(
// 		cashflow.CaptureMessage(sentry.CaptureMessage),
// 		cashflow.CaptureExeption(sentry.CaptureException),
// 		cashflow.FileUpload(uploader.UploadParams{
// 			Tags:         []string{"cashflow"},
// 			Folder:       "uhomestay/cashflows",
// 			ResourceType: "raw",
// 		}, cld.Upload.Upload),
// 		cashflowRepository,
// 	)

// 	duesDeps := dues.NewDeps(
// 		dues.CaptureMessage(sentry.CaptureMessage),
// 		dues.CaptureExeption(sentry.CaptureException),
// 		dues.FileUpload(uploader.UploadParams{
// 			Tags:         []string{"dues"},
// 			Folder:       "uhomestay/dues",
// 			ResourceType: "raw",
// 		}, cld.Upload.Upload),
// 		duesRepository,
// 		memberDuesRepository,
// 		memberRepository,
// 		cashflowRepository,
// 	)

// 	imageDeps := image.NewDeps(
// 		image.CaptureMessage(sentry.CaptureMessage),
// 		image.CaptureExeption(sentry.CaptureException),
// 		image.FileUpload(uploader.UploadParams{
// 			Tags:         []string{"image"},
// 			Folder:       "uhomestay/images-gallery",
// 			ResourceType: "raw",
// 		}, cld.Upload.Upload),
// 		imageRepository,
// 	)

// 	homestayDeps := homestay.NewDeps(
// 		homestay.CaptureMessage(sentry.CaptureMessage),
// 		homestay.CaptureExeption(sentry.CaptureException),
// 		homestay.FileUpload(uploader.UploadParams{
// 			Tags:         []string{"homestay"},
// 			Folder:       "uhomestay/homestay",
// 			ResourceType: "raw",
// 		}, cld.Upload.Upload),
// 		homestayImageRepository,
// 		memberHomestayRepository,
// 		memberRepository,
// 	)

// 	dashboardDeps := dashboard.NewDeps(
// 		dashboard.CaptureMessage(sentry.CaptureMessage),
// 		dashboard.CaptureExeption(sentry.CaptureException),
// 		historyDeps,
// 		imageDeps,
// 		homestayDeps,
// 		documentDeps,
// 		articleDeps,
// 		cashflowDeps,
// 		duesDeps,
// 		userDeps,
// 	)

// 	restApi := handler.NewRestApi(
// 		buildDate,
// 		commitHash,
// 		conf,
// 		posgrePool,
// 		dashboardDeps,
// 	)

// 	restApi.RestApiHandler()
// }


package main

import (
	"os"

	"github.com/gofiber/fiber/v2"
)

func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = ":3000"
	} else {
		port = ":" + port
	}

	return port
}

func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Hello, Railway!",
		})
	})

	app.Listen(getPort())
}