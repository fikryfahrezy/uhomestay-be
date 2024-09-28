package homestay

import (
	"context"
	"io"

	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/user"
	"github.com/cloudinary/cloudinary-go/api/uploader"
)

type (
	FileUploader func(filename string, file io.Reader) (string, error)
)

type HomestayDeps struct {
	Upload                   FileUploader
	HomestayImageRepository  *HomestayImageRepository
	MemberHomestayRepository *MemberHomestayRepository
	MemberRepository         *user.MemberRepository
}

func NewDeps(
	upload FileUploader,
	homestayImageRepository *HomestayImageRepository,
	memberHomestayRepository *MemberHomestayRepository,
	memberRepository *user.MemberRepository,
) *HomestayDeps {
	return &HomestayDeps{
		Upload:                   upload,
		HomestayImageRepository:  homestayImageRepository,
		MemberHomestayRepository: memberHomestayRepository,
		MemberRepository:         memberRepository,
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
