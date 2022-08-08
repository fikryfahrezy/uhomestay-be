package user

import (
	"database/sql"
	"time"

	pgtypeuuid "github.com/jackc/pgtype/ext/gofrs-uuid"
)

type MemberModel struct {
	IsAdmin       bool
	IsApproved    bool
	Name          string
	OtherPhone    string
	WaPhone       string
	ProfilePicUrl string
	IdCardUrl     string
	Username      string
	Password      string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     sql.NullTime
	Id            pgtypeuuid.UUID
}

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

type HomestayImageModel struct {
	Id               uint64
	Name             string
	AlphnumName      string
	Url              string
	CreatedAt        time.Time
	MemberHomestayId sql.NullInt64
	DeletedAt        sql.NullTime
}
