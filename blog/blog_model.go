package blog

import (
	"time"
)

type BlogModel struct {
	Id           string                 `bson:"_id,omitempty"`
	Title        string                 `bson:"title,omitempty"`
	ShortDesc    string                 `bson:"short_desc,omitempty"`
	ThumbnailUrl string                 `bson:"thumbnail_url,omitempty"`
	Slug         string                 `bson:"slug,omitempty"`
	CreatedAt    time.Time              `bson:"created_at,omitempty"`
	UpdatedAt    time.Time              `bson:"updated_at,omitempty"`
	DeletedAt    time.Time              `bson:"deleted_at,omitempty"`
	Content      map[string]interface{} `bson:"content,omitempty"`
}
