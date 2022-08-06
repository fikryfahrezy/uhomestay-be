package homestay

import (
	"database/sql"
	"time"
)

type MemberHomestayModel struct {
	Id           uint64
	Name         string
	Address      string
	Latitude     string
	Longitude    string
	ThumbnailUrl string
	MemberId     string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    sql.NullTime
}
