package cashflow

import (
	"context"
	"io"

	"github.com/cloudinary/cloudinary-go/api/uploader"
	"github.com/getsentry/sentry-go"
)

type (
	FileUploader      func(filename string, file io.Reader) (string, error)
	ExceptionCapturer func(exception error)
	MessageCapturer   func(message string)
)

type CashflowDeps struct {
	CaptureMessage     MessageCapturer
	CaptureExeption    ExceptionCapturer
	Upload             FileUploader
	CashflowRepository *CashflowRepository
}

func NewDeps(
	captureMessage MessageCapturer,
	captureExeption ExceptionCapturer,
	upload FileUploader,
	cashflowRepository *CashflowRepository,
) *CashflowDeps {
	return &CashflowDeps{
		CaptureMessage:     captureMessage,
		CaptureExeption:    captureExeption,
		Upload:             upload,
		CashflowRepository: cashflowRepository,
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
