package cashflow

import (
	"context"
	"io"

	"github.com/cloudinary/cloudinary-go/api/uploader"
)

type Uploader = func(filename string, file io.Reader) (string, error)

type CashflowDeps struct {
	Upload             Uploader
	CashflowRepository *CashflowRepository
}

func NewDeps(upload Uploader, cashflowRepository *CashflowRepository) *CashflowDeps {
	return &CashflowDeps{
		Upload:             upload,
		CashflowRepository: cashflowRepository,
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