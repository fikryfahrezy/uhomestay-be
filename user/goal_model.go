package user

import (
	"time"
)

type GoalModel struct {
	Id          uint64
	OrgPeriodId uint64
	CreatedAt   time.Time
	Vision      map[string]interface{}
	Mission     map[string]interface{}
}
