package dues

import (
	"context"
	"io"

	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/cashflow"
	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/user"
	"github.com/cloudinary/cloudinary-go/api/uploader"
)

type Uploader = func(filename string, file io.Reader) (string, error)

type DuesDeps struct {
	Upload               Uploader
	DuesRepository       *DuesRepository
	MemberDuesRepository *MemberDuesRepository
	MemberRepository     *user.MemberRepository
	CashflowRepository   *cashflow.CashflowRepository
}

func NewDeps(
	upload Uploader,
	duesRepository *DuesRepository,
	memberDuesRepository *MemberDuesRepository,
	memberRepository *user.MemberRepository,
	cashflowRepository *cashflow.CashflowRepository,
) *DuesDeps {
	return &DuesDeps{
		Upload:               upload,
		DuesRepository:       duesRepository,
		MemberDuesRepository: memberDuesRepository,
		MemberRepository:     memberRepository,
		CashflowRepository:   cashflowRepository,
	}
}

func FileUploader(uploadParams uploader.UploadParams, upload func(ctx context.Context, file interface{}, uploadParams uploader.UploadParams) (*uploader.UploadResult, error)) func(filename string, file io.Reader) (url string, err error) {
	return func(filename string, file io.Reader) (url string, err error) {
		ctx := context.Background()
		uploadParams.PublicID = filename

		resp, err := upload(ctx, file, uploadParams)
		if err != nil {
			return "", err
		}

		return resp.SecureURL, nil
	}
}
