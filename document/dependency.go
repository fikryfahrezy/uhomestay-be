package document

import (
	"context"
	"io"

	"github.com/cloudinary/cloudinary-go/api/uploader"
)

type (
	FileUploader func(filename string, file io.Reader) (string, error)
)

type DocumentDeps struct {
	Upload             FileUploader
	DocumentRepository *DocumentRepository
}

func NewDeps(
	upload FileUploader,
	documentRepository *DocumentRepository,
) *DocumentDeps {
	return &DocumentDeps{
		Upload:             upload,
		DocumentRepository: documentRepository,
	}
}

func FileUpload(uploadParams uploader.UploadParams, upload func(ctx context.Context, file interface{}, uploadParams uploader.UploadParams) (*uploader.UploadResult, error)) FileUploader {
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
