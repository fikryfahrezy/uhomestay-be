package user

import (
	"database/sql"
	"time"
)

type OrgStructureModel struct {
	Id            uint64
	PositionId    uint64
	OrgPeriodId   uint64
	PositionLevel int16
	PositionName  string
	MemberId      string
	CreatedAt     time.Time
	DeletedAt     sql.NullTime
}
