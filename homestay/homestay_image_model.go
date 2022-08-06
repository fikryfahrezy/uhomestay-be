package homestay

import (
	"database/sql"
	"time"
)

type HomestayImageModel struct {
	Id               uint64
	Name             string
	AlphnumName      string
	Url              string
	CreatedAt        time.Time
	MemberHomestayId sql.NullInt64
	DeletedAt        sql.NullTime
}
