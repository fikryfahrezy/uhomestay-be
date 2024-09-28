package image

import (
	"context"
	"io"

	"github.com/cloudinary/cloudinary-go/api/uploader"
)

type (
	FileUploader func(filename string, file io.Reader) (string, error)
)

type ImageDeps struct {
	Upload          FileUploader
	ImageRepository *ImageRepository
}

func NewDeps(
	upload FileUploader,
	imageRepository *ImageRepository,
) *ImageDeps {
	return &ImageDeps{
		Upload:          upload,
		ImageRepository: imageRepository,
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
