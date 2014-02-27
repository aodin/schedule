package schedule

import (
	"fmt"
	"time"
)

type Status struct {
	Error error
	Start time.Time
	End   time.Time
}

func (s Status) String() string {
	// TODO Are the casts needed?
	elapsed := float64(s.End.Sub(s.Start).Nanoseconds()) / float64(time.Millisecond)
	if s.Error == nil {
		return fmt.Sprintf("OK (%.3f ms)", elapsed)
	}
	return fmt.Sprintf("ERROR: %s (%.3f ms)", s.Error.Error(), elapsed)
}
