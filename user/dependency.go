package user

import (
	"context"
	"embed"
	"io"
	"path/filepath"

	"github.com/cloudinary/cloudinary-go/api/uploader"
	"github.com/getsentry/sentry-go"
)

type (
	FileUploader      func(filename string, file io.Reader) (string, error)
	ExceptionCapturer func(exception error)
	MessageCapturer   func(message string)
)

type UserDeps struct {
	JwtKey                 []byte
	JwtIssuerUrl           string
	Argon2Salt             string
	JwtAudiences           []string
	CaptureMessage         MessageCapturer
	CaptureExeption        ExceptionCapturer
	Upload                 FileUploader
	Tmpl                   embed.FS
	MemberRepository       *MemberRepository
	PositionRepository     *PositionRepository
	OrgStructureRepository *OrgStructureRepository
	OrgPeriodRepository    *OrgPeriodRepository
	GoalRepository         *GoalRepository
}

func NewDeps(
	jwtKey []byte,
	jwtIssuerUrl string,
	argon2Salt string,
	jwtAudiences []string,
	captureMessage MessageCapturer,
	captureExeption ExceptionCapturer,
	upload FileUploader,
	tmpl embed.FS,
	memberRepository *MemberRepository,
	positionRepository *PositionRepository,
	orgStructureRepository *OrgStructureRepository,
	orgPeriodRepository *OrgPeriodRepository,
	goalRepository *GoalRepository,
) *UserDeps {
	return &UserDeps{
		JwtKey:                 jwtKey,
		JwtIssuerUrl:           jwtIssuerUrl,
		Argon2Salt:             argon2Salt,
		CaptureMessage:         captureMessage,
		CaptureExeption:        captureExeption,
		JwtAudiences:           jwtAudiences,
		Upload:                 upload,
		Tmpl:                   tmpl,
		MemberRepository:       memberRepository,
		PositionRepository:     positionRepository,
		OrgStructureRepository: orgStructureRepository,
		OrgPeriodRepository:    orgPeriodRepository,
		GoalRepository:         goalRepository,
	}
}

func FileUpload(uploadParams uploader.UploadParams, upload func(ctx context.Context, file interface{}, uploadParams uploader.UploadParams) (*uploader.UploadResult, error)) FileUploader {
	return func(filename string, file io.Reader) (url string, err error) {
		filename = filename[:len(filename)-len(filepath.Ext(filename))]
		ctx := context.Background()
		uploadParams.PublicID = filename

		resp, err := upload(ctx, file, uploadParams)
		if err != nil {
			return "", err
		}

		return resp.SecureURL, nil
	}
}

func CaptureExeption(capture func(exception error) *sentry.EventID) ExceptionCapturer {
	return func(exception error) {
		capture(exception)
	}
}

func CaptureMessage(capture func(message string) *sentry.EventID) MessageCapturer {
	return func(message string) {
		capture(message)
	}
}
