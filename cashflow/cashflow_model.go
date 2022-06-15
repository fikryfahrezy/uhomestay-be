package cashflow

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"time"
)

// Ref: Saving enumerated values to a database
// https://stackoverflow.com/a/25374979/12976234
type CashflowType struct {
	String string
}

var (
	Unknown = CashflowType{""}
	Income  = CashflowType{"income"}
	Outcome = CashflowType{"outcome"}
)

func typeFromString(s string) (CashflowType, error) {
	switch s {
	case Income.String:
		return Income, nil
	case Outcome.String:
		return Outcome, nil
	}

	return Unknown, errors.New("unknown type: " + s)
}

func (u *CashflowType) Scan(src interface{}) error {
	if src == nil {
		u.String = ""
		return nil
	}

	s, ok := src.(string)
	if !ok {
		u.String = ""
		return nil
	}

	dc, _ := typeFromString(s)
	u.String = dc.String
	return nil
}

func (u CashflowType) Value() (driver.Value, error) {
	dc, err := typeFromString(u.String)
	if err != nil {
		dc = Income
	}

	return dc.String, nil
}

type CashflowModel struct {
	Id           uint64
	IdrAmount    string
	Note         string
	ProveFileUrl string
	Type         CashflowType
	Date         time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    sql.NullTime
}
