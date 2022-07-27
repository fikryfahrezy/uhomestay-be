package dues

import (
	"database/sql"
	"time"
)

type MemberDuesViewModel struct {
	Id           uint64
	DuesId       uint64
	IdrAmount    string
	ProveFileUrl string
	Status       DuesStatus
	Date         time.Time
	PayDate      sql.NullTime
}

type DuesMemberViewModel struct {
	Id            uint64
	MemberId      string
	Name          string
	ProfilePicUrl string
	Status        DuesStatus
	CreatedAt     time.Time
	PayDate       sql.NullTime
}

type MemberDuesAmtViewModel struct {
	IdrAmount string
	Status    DuesStatus
}
