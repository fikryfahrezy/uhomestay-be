package dues

import (
	"database/sql"
	"database/sql/driver"
	"time"

	"github.com/pkg/errors"
)

// Ref: Saving enumerated values to a database
// https://stackoverflow.com/a/25374979/12976234
type DuesStatus struct {
	String string
}

var (
	Unknown = DuesStatus{""}
	Unpaid  = DuesStatus{"unpaid"}
	Waiting = DuesStatus{"waiting"}
	Paid    = DuesStatus{"paid"}
)

func typeFromString(s string) (DuesStatus, error) {
	switch s {
	case Unpaid.String:
		return Unpaid, nil
	case Waiting.String:
		return Waiting, nil
	case Paid.String:
		return Paid, nil
	}

	return Unknown, errors.New("unknown type: " + s)
}

func (u *DuesStatus) Scan(src interface{}) error {
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

func (u DuesStatus) Value() (driver.Value, error) {
	dc, err := typeFromString(u.String)
	if err != nil {
		dc = Unpaid
	}

	return dc.String, nil
}

type MemberDuesModel struct {
	Id           uint64
	DuesId       uint64
	ProveFileUrl string
	MemberId     string
	Status       DuesStatus
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    sql.NullTime
}
