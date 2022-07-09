package history

import "github.com/getsentry/sentry-go"

type (
	ExceptionCapturer func(exception error)
	MessageCapturer   func(message string)
)

type HistoryDeps struct {
	CaptureMessage    MessageCapturer
	CaptureExeption   ExceptionCapturer
	HistoryRepository *HistoryRepository
}

func NewDeps(
	captureMessage MessageCapturer,
	captureExeption ExceptionCapturer,
	historyRepository *HistoryRepository,
) *HistoryDeps {
	return &HistoryDeps{
		CaptureMessage:    captureMessage,
		CaptureExeption:   captureExeption,
		HistoryRepository: historyRepository,
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
