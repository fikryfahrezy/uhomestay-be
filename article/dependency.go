package article

import (
	"context"
	"io"

	"github.com/cloudinary/cloudinary-go/api/uploader"
	"github.com/getsentry/sentry-go"
)

type (
	FileUploader      func(filename string, file io.Reader) (string, string, error)
	FileMover         func(from, to string) (string, error)
	ExceptionCapturer func(exception error)
	MessageCapturer   func(message string)
)

type ArticleDeps struct {
	ImgCldTmpFolder   string
	ImgClgFolder      string
	CaptureMessage    MessageCapturer
	CaptureExeption   ExceptionCapturer
	MoveFile          FileMover
	Upload            FileUploader
	ArticleRepository *ArticleRepository
}

func NewDeps(
	imgClgFolder string,
	imgCldTmpFolder string,
	captureMessage MessageCapturer,
	captureExeption ExceptionCapturer,
	moveFile FileMover,
	upload FileUploader,
	articleRepository *ArticleRepository,
) *ArticleDeps {
	return &ArticleDeps{
		ImgClgFolder:      imgClgFolder,
		ImgCldTmpFolder:   imgCldTmpFolder,
		CaptureMessage:    captureMessage,
		CaptureExeption:   captureExeption,
		MoveFile:          moveFile,
		Upload:            upload,
		ArticleRepository: articleRepository,
	}
}

func FileUpload(uploadParams uploader.UploadParams, upload func(ctx context.Context, file interface{}, uploadParams uploader.UploadParams) (*uploader.UploadResult, error)) FileUploader {
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

func FileMove(move func(ctx context.Context, params uploader.RenameParams) (*uploader.RenameResult, error)) FileMover {
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
