package user

import (
	"context"
	"embed"
	"io"
	"path/filepath"

	"github.com/cloudinary/cloudinary-go/api/uploader"
)

type Uploader = func(filename string, file io.Reader) (string, error)

type UserDeps struct {
	JwtKey                 []byte
	JwtIssuerUrl           string
	Argon2Salt             string
	JwtAudiences           []string
	Upload                 Uploader
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
	upload Uploader,
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

func FileUploader(uploadParams uploader.UploadParams, upload func(ctx context.Context, file interface{}, uploadParams uploader.UploadParams) (*uploader.UploadResult, error)) func(filename string, file io.Reader) (url string, err error) {
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
