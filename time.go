package schedule

import (
	"time"
)

// TickAt is a wrappers for time channels. If the time has already occured,
// the current time will be sent immediately
func TickAt(t time.Time) <-chan time.Time {
	return tickAt(defaultNow, t)
}

// defaultNow allows composition of the now() function for testing.
func defaultNow() time.Time {
	return time.Now()
}

func tickAt(now func() time.Time, t time.Time) <-chan time.Time {
	delta := t.Sub(now())
	if delta < 0 {
		delta = time.Duration(0)
	}
	return time.NewTimer(delta).C
}
