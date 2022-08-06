package dashboard

import (
	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/article"
	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/cashflow"
	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/document"
	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/dues"
	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/history"
	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/homestay"
	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/image"
	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/user"
	"github.com/getsentry/sentry-go"
)

type (
	ExceptionCapturer func(exception error)
	MessageCapturer   func(message string)
)

type DashboardDeps struct {
	CaptureMessage  MessageCapturer
	CaptureExeption ExceptionCapturer
	*history.HistoryDeps
	*image.ImageDeps
	*homestay.HomestayDeps
	*document.DocumentDeps
	*article.ArticleDeps
	*cashflow.CashflowDeps
	*dues.DuesDeps
	*user.UserDeps
}

func NewDeps(
	captureMessage MessageCapturer,
	captureExeption ExceptionCapturer,
	historyDeps *history.HistoryDeps,
	imageDeps *image.ImageDeps,
	homestayDeps *homestay.HomestayDeps,
	documentDeps *document.DocumentDeps,
	articleDeps *article.ArticleDeps,
	cashflowDeps *cashflow.CashflowDeps,
	duesDeps *dues.DuesDeps,
	userDeps *user.UserDeps,
) *DashboardDeps {
	return &DashboardDeps{
		CaptureMessage:  captureMessage,
		CaptureExeption: captureExeption,
		HistoryDeps:     historyDeps,
		ImageDeps:       imageDeps,
		HomestayDeps:    homestayDeps,
		DocumentDeps:    documentDeps,
		ArticleDeps:     articleDeps,
		CashflowDeps:    cashflowDeps,
		DuesDeps:        duesDeps,
		UserDeps:        userDeps,
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
