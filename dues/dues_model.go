package dues

import (
	"database/sql"
	"time"
)

type DuesModel struct {
	Id        uint64
	IdrAmount string
	Date      time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt sql.NullTime
}
