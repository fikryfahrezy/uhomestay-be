package blog

import (
	"context"
	"io"

	"github.com/cloudinary/cloudinary-go/api/uploader"
)

type (
	Uploader = func(filename string, file io.Reader) (string, string, error)
	Mover    = func(from, to string) (string, error)
)

type BlogDeps struct {
	ImgCldTmpFolder string
	ImgClgFolder    string
	MoveFile        Mover
	Upload          Uploader
	BlogRepository  *BlogRepository
}

func NewDeps(
	imgClgFolder string,
	imgCldTmpFolder string,
	moveFile Mover,
	upload Uploader,
	blogRepository *BlogRepository,
) *BlogDeps {
	return &BlogDeps{
		ImgClgFolder:    imgClgFolder,
		ImgCldTmpFolder: imgCldTmpFolder,
		MoveFile:        moveFile,
		Upload:          upload,
		BlogRepository:  blogRepository,
	}
}

func FileUploader(uploadParams uploader.UploadParams, upload func(ctx context.Context, file interface{}, uploadParams uploader.UploadParams) (*uploader.UploadResult, error)) func(filename string, file io.Reader) (url string, id string, err error) {
	return func(filename string, file io.Reader) (url string, id string, err error) {
		ctx := context.Background()
		uploadParams.PublicID = filename

		resp, err := upload(ctx, file, uploadParams)
		if err != nil {
			return "", "", err
		}

		return resp.SecureURL, resp.PublicID, nil
	}
}

func FileMover(move func(ctx context.Context, params uploader.RenameParams) (*uploader.RenameResult, error)) func(from, to string) (url string, err error) {
	return func(from, to string) (url string, err error) {
		ctx := context.Background()

		resp, err := move(ctx, uploader.RenameParams{
			FromPublicID: from,
			ToPublicID:   to,
		})
		if err != nil {
			return "", err
		}

		return resp.SecureURL, nil
	}
}
