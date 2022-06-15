package timediff

import (
	"time"
)

// Return `true` if `currDate` is lower than `cmpDate`
// Return `false` otherwise
func Compare(currDate, cmpDate time.Time) bool {
	year, month, _ := currDate.Date()
	currY, currM, _ := cmpDate.Date()
	prevY, prevM, _ := time.Date(currY, currM, 1, 0, 0, 0, 0, time.UTC).Date()

	im, ipm := int(month), int(prevM)
	if (year-prevY) < 0 || (year == prevY && im < ipm) {
		return true
	}

	return false
}

func IsPast(currDate time.Time) bool {
	compared := Compare(currDate, time.Now())
	return compared
}
