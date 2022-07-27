package document

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"time"
)

// Ref: Saving enumerated values to a database
// https://stackoverflow.com/a/25374979/12976234
type DocType struct {
	String string
}

var (
	Unknown  = DocType{""}
	Dir      = DocType{"dir"}
	Filetype = DocType{"file"}
)

func typeFromString(s string) (DocType, error) {
	switch s {
	case Dir.String:
		return Dir, nil
	case Filetype.String:
		return Filetype, nil
	}

	return Unknown, errors.New("unknown type: " + s)
}

func (u *DocType) Scan(src interface{}) error {
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

func (u DocType) Value() (driver.Value, error) {
	dc, err := typeFromString(u.String)
	if err != nil {
		dc = Filetype
	}

	return dc.String, nil
}

type DocumentModel struct {
	IsPrivate   bool
	Id          uint64
	DirId       uint64
	Name        string
	AlphnumName string
	Url         string
	Type        DocType
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   sql.NullTime
}
