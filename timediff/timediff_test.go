package timediff_test

import (
	"testing"
	"time"

	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/timediff"
)

func TestCompare(t *testing.T) {
	testCases := []struct {
		res   bool
		name  string
		start time.Time
		end   time.Time
	}{
		{
			res:   true,
			name:  "start 2009-11-10 10:23:0:0:0 UTC | end 2010-01-10 10:23:0:0:0 UTC",
			start: time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
			end:   time.Date(2010, time.January, 10, 23, 0, 0, 0, time.UTC),
		},
		{
			res:   true,
			name:  "start 2009-11-10 10:23:0:0:0 UTC | end 2009-12-10 10:23:0:0:0 UTC",
			start: time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
			end:   time.Date(2009, time.December, 10, 23, 0, 0, 0, time.UTC),
		},
		{
			res:   true,
			name:  "start 2008-12-10 10:23:0:0:0 UTC | end 2009-11-10 10:23:0:0:0 UTC",
			start: time.Date(2008, time.December, 10, 23, 0, 0, 0, time.UTC),
			end:   time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
		},
		{
			res:   false,
			name:  "start 2009-12-10 10:23:0:0:0 UTC | end 2009-11-10 10:23:0:0:0 UTC",
			start: time.Date(2009, time.December, 10, 23, 0, 0, 0, time.UTC),
			end:   time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
		},
		{
			res:   false,
			name:  "start 2011-12-10 10:23:0:0:0 UTC | end 2009-11-10 10:23:0:0:0 UTC",
			start: time.Date(2011, time.December, 10, 23, 0, 0, 0, time.UTC),
			end:   time.Date(2009, time.December, 10, 23, 0, 0, 0, time.UTC),
		},
		{
			res:   false,
			name:  "start 2009-12-10 11:23:0:0:0 UTC | end 2009-12-10 10:23:0:0:0 UTC",
			start: time.Date(2009, time.December, 11, 23, 0, 0, 0, time.UTC),
			end:   time.Date(2009, time.December, 10, 23, 0, 0, 0, time.UTC),
		},
		{
			res:   false,
			name:  "start 2009-12-10 11:23:0:0:1 UTC | end 2009-12-10 11:23:0:0:0 UTC",
			start: time.Date(2009, time.December, 11, 23, 0, 0, 1, time.UTC),
			end:   time.Date(2009, time.December, 11, 23, 0, 0, 0, time.UTC),
		},
	}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			res := timediff.Compare(c.start, c.end)
			if c.res != res {
				t.Fatalf("start: %v, end: %v, resulting: %v, expect: %v", c.start.String(), c.end.String(), res, c.res)
			}
		})
	}
}

func TestIsPast(t *testing.T) {
	testCases := []struct {
		res  bool
		name string
		in   time.Time
	}{
		{
			res:  true,
			name: "in 2009-12-10 10:23:0:0:0 UTC",
			in:   time.Date(2009, time.December, 10, 23, 0, 0, 0, time.UTC),
		},
		{
			res:  true,
			name: "in 2011-12-10 10:23:0:0:0 UTC",
			in:   time.Date(2011, time.December, 10, 23, 0, 0, 0, time.UTC),
		},
		{
			res:  true,
			name: "in 2009-12-10 11:23:0:0:0 UTC",
			in:   time.Date(2009, time.December, 11, 23, 0, 0, 0, time.UTC),
		},
		{
			res:  true,
			name: "in 2009-12-10 11:23:0:0:1 UTC",
			in:   time.Date(2009, time.December, 11, 23, 0, 0, 1, time.UTC),
		},
		{
			res:  true,
			name: "in 2009-12-10 11:23:0:0:1 UTC",
			in:   time.Date(2009, time.December, 11, 23, 0, 0, 1, time.UTC),
		},
		{
			res:  false,
			name: "in time now",
			in:   time.Now(),
		},
		{
			res:  false,
			name: "in time now + 1 hour",
			in:   time.Now().Add(time.Hour),
		},
		{
			res:  false,
			name: "in time now + 1 day",
			in:   time.Now().Add(24 * time.Hour),
		},
		{
			res:  false,
			name: "in time now + 1 month",
			in:   time.Now().Add(24 * time.Hour * 31),
		},
		{
			res:  false,
			name: "in time now + 1 year",
			in:   time.Now().Add(24 * time.Hour * 31 * 12),
		},
		{
			res:  false,
			name: "in time now + 2 year",
			in:   time.Now().Add((24 * time.Hour * 31 * 12) * 2),
		},
	}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			res := timediff.IsPast(c.in)
			if c.res != res {
				t.Fatalf("in: %v, resulting: %v, expect: %v", c.in.String(), res, c.res)
			}
		})
	}
}
