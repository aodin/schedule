package schedule

import (
	"time"
)

type Daytime struct {
	day   time.Weekday
	clock Clock
}

// The location of a Daytime is determined by its Clock
func (d Daytime) Location() *time.Location {
	return d.clock.loc
}

// Return the next occurrence of this day and clock
func (d Daytime) Next() time.Time {
	return d.next(defaultNow)
}

func (d Daytime) next(now func() time.Time) time.Time {
	n := now()
	// TODO If this clock has already occurred today, add a day
	return n
}

func FromDayAndClock(day time.Weekday, clock Clock) Daytime {
	return Daytime{day, clock}
}
