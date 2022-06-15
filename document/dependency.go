package document

import (
	"context"
	"io"

	"github.com/cloudinary/cloudinary-go/api/uploader"
)

type Uploader func(filename string, file io.Reader) (string, error)

type DocumentDeps struct {
	Upload             Uploader
	DocumentRepository *DocumentRepository
}

func NewDeps(upload Uploader, documentRepository *DocumentRepository) *DocumentDeps {
	return &DocumentDeps{
		Upload:             upload,
		DocumentRepository: documentRepository,
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
