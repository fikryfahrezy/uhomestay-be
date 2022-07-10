package pagination_test

import (
	"testing"
	"time"

	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/pagination"
)

func TestDecodeSIDCursor(t *testing.T) {
	testCases := []struct {
		name      string
		inString  string
		outString string
		cursor    string
		inTime    time.Time
		outTime   time.Time
	}{
		{
			name:      "Decode string cursor success",
			inString:  "test",
			outString: "test",
			cursor:    pagination.EncodeSIDCursor("test", time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)),
			inTime:    time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
			outTime:   time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
		},
		{
			name:      "Decode string cursor success",
			outString: "",
			cursor:    "",
			inTime:    time.Time{},
			outTime:   time.Time{},
		},
		{
			name:      "Decode string cursor success, but resulting empty string",
			inString:  "",
			outString: "",
			cursor:    "74657374",
			inTime:    time.Time{},
			outTime:   time.Time{},
		},
		{
			name:      "Decode base64 sucess, but string cursor not as expected",
			inString:  "",
			outString: "",
			cursor:    "dGVzdA==",
			inTime:    time.Time{},
			outTime:   time.Time{},
		},
		{
			name:      "Decode base64 sucess, but invalid time",
			inString:  "",
			outString: "",
			cursor:    "dGVzdCx0ZXN0",
			inTime:    time.Time{},
			outTime:   time.Time{},
		},
	}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			outS, outT, err := pagination.DecodeSIDCursor(c.cursor)
			if c.inString != outS {
				t.Fatalf("string, input: %v, expect: %v, result: %v, err: %v", c.inString, c.outString, outS, err)
			}

			if c.inTime.String() != outT.String() {
				t.Fatalf("string, input: %v, expect: %v, result: %v, err: %v", c.inTime.String(), c.outTime.String(), outT.String(), err)
			}
		})
	}
}
