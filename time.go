package schedule

import (
	"time"
)

// Wrappers around time channels
// If the time has already occured, the current time will be sent immediately
// TODO Return the timer so that it can be stopped
// The returned tick should be the scheduled time
func TickAt(t time.Time) <-chan time.Time {
	return tickAt(defaultNow, t)
}

// Allow composition of the now() function for testing
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
