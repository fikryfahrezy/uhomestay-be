package user

import (
	"database/sql"
	"time"
)

type PositionModel struct {
	Id        uint64
	Name      string
	Level     int16
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt sql.NullTime
}
