package schedule

import (
	"time"
)

// Daytime is combination of day of the week and time of day.
type Daytime struct {
	day   time.Weekday
	clock Clock
}

// Location returns the location of a Daytime. Location is set by its Clock.
func (d Daytime) Location() *time.Location {
	return d.clock.loc
}

// Next returns the next occurrence of this day and clock.
func (d Daytime) Next() time.Time {
	return d.next(defaultNow)
}

func (d Daytime) next(now func() time.Time) time.Time {
	n := now()
	// TODO If this clock has already occurred today, add a day
	return n
}

// FromDayAndClock creates a Daytime from a given day of the week and time
// of day.
func FromDayAndClock(day time.Weekday, clock Clock) Daytime {
	return Daytime{day, clock}
}
