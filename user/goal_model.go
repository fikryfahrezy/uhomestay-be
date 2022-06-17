package user

import (
	"time"
)

type GoalModel struct {
	Id          uint64
	OrgPeriodId uint64
	VisionText  string
	MissionText string
	CreatedAt   time.Time
	Vision      map[string]interface{}
	Mission     map[string]interface{}
}
