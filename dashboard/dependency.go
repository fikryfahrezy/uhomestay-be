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
)

type DashboardDeps struct {
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
		HistoryDeps:  historyDeps,
		ImageDeps:    imageDeps,
		HomestayDeps: homestayDeps,
		DocumentDeps: documentDeps,
		ArticleDeps:  articleDeps,
		CashflowDeps: cashflowDeps,
		DuesDeps:     duesDeps,
		UserDeps:     userDeps,
	}
}
