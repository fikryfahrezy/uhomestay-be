package image

import (
	"database/sql"
	"time"
)

type ImageModel struct {
	Id          uint64
	Name        string
	AlphnumName string
	Url         string
	Description string
	CreatedAt   time.Time
	DeletedAt   sql.NullTime
}
