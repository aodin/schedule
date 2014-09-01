package schedule

import (
	"fmt"
	"time"
)

// Status records the start and end time of a task. It will include the
// task's error message if one occurred.
type Status struct {
	Error error
	Start time.Time
	End   time.Time
}

// String returns a basic string with the task's elapsed time and error
// message if one occurred.
func (s Status) String() string {
	// TODO Are the casts needed?
	elapsed := float64(s.End.Sub(s.Start).Nanoseconds()) / float64(time.Millisecond)
	if s.Error == nil {
		return fmt.Sprintf("OK (%.3f ms)", elapsed)
	}
	return fmt.Sprintf("ERROR: %s (%.3f ms)", s.Error.Error(), elapsed)
}
