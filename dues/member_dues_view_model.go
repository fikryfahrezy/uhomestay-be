package dues

import (
	"time"
)

type MemberDuesViewModel struct {
	Id           uint64
	DuesId       uint64
	IdrAmount    string
	ProveFileUrl string
	Status       DuesStatus
	Date         time.Time
}

type DuesMemberViewModel struct {
	Id            uint64
	MemberId      string
	Name          string
	ProfilePicUrl string
	Status        DuesStatus
	CreatedAt     time.Time
}
