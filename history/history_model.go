package history

import (
	"time"
)

type HistoryModel struct {
	Id        string                 `bson:"_id,omitempty"`
	CreatedAt time.Time              `bson:"created_at,omitempty"`
	Content   map[string]interface{} `bson:"content,omitempty"`
}
