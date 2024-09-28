package main

import (
	"context"
	"embed"
	"log"

	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/article"
	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/cashflow"
	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/config"
	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/dashboard"
	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/document"
	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/dues"
	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/handler"
	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/history"
	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/homestay"
	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/image"
	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/user"

	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/uploader"
	"github.com/jackc/pgtype"
	pgtypeuuid "github.com/jackc/pgtype/ext/gofrs-uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

//go:embed tmpl/*
var tmpl embed.FS

var (
	buildDate  string = "N/A"
	commitHash string = "N/A"
)

func main() {
	conf := config.LoadConfig()

	posgreConfig, err := pgxpool.ParseConfig(conf.PostgreUrl)
	if err != nil {
		log.Fatalf("fail parse posgre database config: %s", err)
	}
	posgreConfig.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		conn.ConnInfo().RegisterDataType(pgtype.DataType{
			Value: &pgtypeuuid.UUID{},
			Name:  "uuid",
			OID:   pgtype.UUIDOID,
		})
		return nil
	}

	posgrePool, err := pgxpool.ConnectConfig(context.Background(), posgreConfig)
	if err != nil {
		log.Fatalf("fail connect to postgre: %s", err)
	}
	defer posgrePool.Close()

	if _, err := posgrePool.Exec(context.Background(), "SELECT 1;"); err != nil {
		log.Fatalf("fail to ping postgre: %s", err)
		return
	}

	cld, err := cloudinary.NewFromURL(conf.CloudinaryUrl)
	if err != nil {
		log.Fatalf("cloudinary.NewFromURL: %s", err)
	}

	memberRepository := user.NewMemberRepository(posgrePool)
	positionRepository := user.NewPositionRepository(posgrePool)
	orgRepository := user.NewOrgStructureRepository(posgrePool)
	periodRepository := user.NewOrgPeriodRepository(posgrePool)
	goalRepository := user.NewGoalRepository(posgrePool)
	documentRepository := document.NewRepository(posgrePool)
	cashflowRepository := cashflow.NewRepository(posgrePool)
	duesRepository := dues.NewDeusRepository(posgrePool)
	memberDuesRepository := dues.NewMemberDeusRepository(posgrePool)
	imageRepository := image.NewRepository(posgrePool)

	historyRepository := history.NewRepository(
		posgrePool,
	)
	articleRepository := article.NewRepository(
		"imgchc",
		posgrePool,
	)

	memberHomestayRepository := homestay.NewMemberHomestayRepository(
		posgrePool,
	)
	homestayImageRepository := homestay.NewHomestayImageRepository(
		posgrePool,
	)

	userDeps := user.NewDeps(
		conf.JwtKey,
		conf.JwtIssuerUrl,
		conf.Argon2Salt,
		conf.JwtAudiences,
		user.FileUpload(uploader.UploadParams{
			Transformation: "c_crop,g_center/q_auto/f_auto",
			Tags:           []string{"profile"},
			Folder:         "uhomestay/profile",
			ResourceType:   "image",
		}, cld.Upload.Upload),
		tmpl,
		memberRepository,
		positionRepository,
		orgRepository,
		periodRepository,
		goalRepository,
	)

	documentDeps := document.NewDeps(
		document.FileUpload(uploader.UploadParams{
			Tags:         []string{"document"},
			Folder:       "uhomestay/document",
			ResourceType: "raw",
		}, cld.Upload.Upload),
		documentRepository,
	)

	historyDeps := history.NewDeps(
		historyRepository,
	)

	articleImgFolder := "uhomestay/blog-images-tmp"
	articleDeps := article.NewDeps(
		"uhomestay/blog-images",
		articleImgFolder,
		article.FileMove(cld.Upload.Rename),
		article.FileUpload(uploader.UploadParams{
			Tags:         []string{"blogs"},
			Folder:       articleImgFolder,
			ResourceType: "raw",
		}, cld.Upload.Upload),
		articleRepository,
	)

	cashflowDeps := cashflow.NewDeps(
		cashflow.FileUpload(uploader.UploadParams{
			Tags:         []string{"cashflow"},
			Folder:       "uhomestay/cashflows",
			ResourceType: "raw",
		}, cld.Upload.Upload),
		cashflowRepository,
	)

	duesDeps := dues.NewDeps(
		dues.FileUpload(uploader.UploadParams{
			Tags:         []string{"dues"},
			Folder:       "uhomestay/dues",
			ResourceType: "raw",
		}, cld.Upload.Upload),
		duesRepository,
		memberDuesRepository,
		memberRepository,
		cashflowRepository,
	)

	imageDeps := image.NewDeps(
		image.FileUpload(uploader.UploadParams{
			Tags:         []string{"image"},
			Folder:       "uhomestay/images-gallery",
			ResourceType: "raw",
		}, cld.Upload.Upload),
		imageRepository,
	)

	homestayDeps := homestay.NewDeps(
		homestay.FileUpload(uploader.UploadParams{
			Tags:         []string{"homestay"},
			Folder:       "uhomestay/homestay",
			ResourceType: "raw",
		}, cld.Upload.Upload),
		homestayImageRepository,
		memberHomestayRepository,
		memberRepository,
	)

	dashboardDeps := dashboard.NewDeps(
		historyDeps,
		imageDeps,
		homestayDeps,
		documentDeps,
		articleDeps,
		cashflowDeps,
		duesDeps,
		userDeps,
	)

	restApi := handler.NewRestApi(
		buildDate,
		commitHash,
		conf,
		posgrePool,
		dashboardDeps,
	)

	restApi.RestApiHandler()
}
