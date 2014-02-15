package schedule

import (
	"time"
)

// Wrappers around time channels
// If the time has already occured, the current time will be sent immediately
// TODO Return the timer so that it can be stopped
// The returned tick should be the scheduled time
func TickAt(t time.Time) <-chan time.Time {
	return tickAt(now, t)
}

// Allow composition of the now() function for testing
func now() time.Time {
	return time.Now()
}

func tickAt(getNow func() time.Time, t time.Time) <-chan time.Time {
	delta := t.Sub(time.Now())
	if delta < 0 {
		delta = time.Duration(0)
	}
	return time.NewTimer(delta).C
}

// Create a channel that will send on the next occurrence of the given clock
// func NextClock(clock Clock) <-chan Time {

// }
