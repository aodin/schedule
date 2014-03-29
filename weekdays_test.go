package schedule

import (
	"testing"
	"time"
)

func TestDaysAway(t *testing.T) {
	setTuesday := func() time.Time {
		return time.Date(2014, time.Month(3), 18, 12, 12, 12, 0, time.Local)
	}
	expectInt(t, daysAway(setTuesday, time.Tuesday), 0)
	expectInt(t, daysAway(setTuesday, time.Wednesday), 1)
	expectInt(t, daysAway(setTuesday, time.Saturday), 4)
	expectInt(t, daysAway(setTuesday, time.Sunday), 5)
	expectInt(t, daysAway(setTuesday, time.Monday), 6)
}
