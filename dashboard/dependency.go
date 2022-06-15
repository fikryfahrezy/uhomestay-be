package dashboard

import (
	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/blog"
	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/cashflow"
	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/document"
	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/dues"
	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/history"
	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/user"
)

type DashboardDeps struct {
	*history.HistoryDeps
	*document.DocumentDeps
	*blog.BlogDeps
	*cashflow.CashflowDeps
	*dues.DuesDeps
	*user.UserDeps
}

func NewDeps(
	historyDeps *history.HistoryDeps,
	documentDeps *document.DocumentDeps,
	blogDeps *blog.BlogDeps,
	cashflowDeps *cashflow.CashflowDeps,
	duesDeps *dues.DuesDeps,
	userDeps *user.UserDeps,
) *DashboardDeps {
	return &DashboardDeps{
		HistoryDeps:  historyDeps,
		DocumentDeps: documentDeps,
		BlogDeps:     blogDeps,
		CashflowDeps: cashflowDeps,
		DuesDeps:     duesDeps,
		UserDeps:     userDeps,
	}
}
