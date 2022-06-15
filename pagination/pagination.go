package pagination

import (
	"encoding/base64"
	"errors"
	"strings"
	"time"
)

// Ref: https://github.com/bxcodec/go-clean-arch/blob/master/article/repository/helper.go
func DecodeSIDCursor(encodedTime string) (s string, t time.Time, err error) {
	byt, err := base64.StdEncoding.DecodeString(encodedTime)
	if err != nil {
		return "", time.Time{}, err
	}

	sbyt := string(byt)
	if sbyt == "" {
		return "", time.Time{}, nil
	}

	arrStr := strings.Split(sbyt, ",")
	if len(arrStr) != 2 {
		err = errors.New("cursor is invalid")
		return "", time.Time{}, err

	}

	timeString := arrStr[1]
	if timeString == "" {
		timeString = "1970-01-01"
	}

	t, err = time.Parse(time.RFC3339, timeString)
	if err != nil {
		return "", time.Time{}, err
	}

	s = arrStr[0]

	return s, t, nil
}

func EncodeSIDCursor(sid string, t time.Time) string {
	timeString := sid + "," + t.Format(time.RFC3339)

	return base64.StdEncoding.EncodeToString([]byte(timeString))
}
