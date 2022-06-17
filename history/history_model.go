package history

import (
	"time"
)

type HistoryModel struct {
	Id          uint64
	ContentText string
	CreatedAt   time.Time
	Content     map[string]interface{}
}
