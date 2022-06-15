package user

import (
	"database/sql"
	"time"
)

type OrgPeriodModel struct {
	IsActive  bool
	Id        uint64
	StartDate time.Time
	EndDate   time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt sql.NullTime
}
