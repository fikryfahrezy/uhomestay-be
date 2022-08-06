package article

import (
	"database/sql"
	"time"
)

type ArticleModel struct {
	Id           uint64
	Title        string
	ShortDesc    string
	ThumbnailUrl string
	ContentText  string
	Slug         string
	Content      map[string]interface{}
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    sql.NullTime
}
