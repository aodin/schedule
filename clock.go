package schedule

import (
	"fmt"
	"sort"
	"time"
)

// TODO Future acceptable time formats
// 5pm
// 5:15pm
// 5:15 pm
// 5:15.00 pm
// 17:15.00000

// Clock stores the seconds since the beginning of the day and a location.
type Clock struct {
	sec int // TODO nano-seconds?
	loc *time.Location
}

// String returns a string with a pretty printed time of day.
func (c Clock) String() string {
	hour, minute, second := c.HMS()
	return fmt.Sprintf("%d:%02d:%02d", hour, minute, second)
}

// Add adds the time.Duration to the current clock.
func (c Clock) Add(d time.Duration) Clock {
	c.sec += int(d / 1e9)
	return c
}

// Before returns a boolean indicating if the clock occurred before the given
// clock.
func (c Clock) Before(other Clock) bool {
	// Clocks may have been built through addition, so mod by seconds in a day
	return (c.sec % (60 * 60 * 24)) < (other.sec % (60 * 60 * 24))
}

// Locations returns the location assigned to the clock.
func (c Clock) Location() *time.Location {
	return c.loc
}

// TODO Return the clock at a different timezone

// HMS returns the hour, minute, and second indicated by this clock.
func (c Clock) HMS() (int, int, int) {
	return c.Hour(), c.Minute(), c.Second()
}

// Hour returns the hour indicated by this clock. If there are more than one
// day's worth of seconds, it will roll over to the next day.
func (c Clock) Hour() int {
	return (c.sec / (60 * 60)) % 24
}

// Minute returns the minute indicated by this clock.
func (c Clock) Minute() int {
	return (c.sec / 60) % 60
}

// Second returns the second indicated by this clock.
func (c Clock) Second() int {
	return c.sec % 60
}

// TotalSeconds returns the total number of seconds indicated by this clock.
func (c Clock) TotalSeconds() int {
	return c.sec
}

// UTC sets the clock's location to UTC. It does not change the clock's time.
func (c Clock) UTC() Clock {
	c.loc = time.UTC
	// TODO This operation does not make much sense - adjust the clock?
	return c
}

// Next returns a `time.Time` for the next occurrence of the clock.
func (c Clock) Next() time.Time {
	// TODO Timezones matter!
	return c.next(defaultNow)
}

func (c Clock) next(now func() time.Time) time.Time {
	n := now().In(c.loc)
	year, month, day := n.Date()
	nxt := c.ToTime(year, month, day)
	if nxt.Before(n) {
		// If nxt has already occured, add a day
		nxt = nxt.Add(24 * time.Hour)
	}
	return nxt
}

// ToTime converts the given clock to a `time.Time` using the given
// year, month, and date.
func (c Clock) ToTime(y int, m time.Month, d int) time.Time {
	return time.Date(y, m, d, c.Hour(), c.Minute(), c.Second(), 0, c.loc)
}

// ClockFromTime creates a Clock from the given `time.Time`.
func ClockFromTime(t time.Time) Clock {
	hr, mm, ss := t.Clock()
	return Clock{(hr*60+mm)*60 + ss, t.Location()}
}

// ClockNow created a Clock of the current local time.
func ClockNow() Clock {
	return clockNowIn(defaultNow, time.Local)
}

// ClockNow created a Clock of the current UTC time.
func ClockNowUTC() Clock {
	return clockNowIn(defaultNow, time.UTC)
}

func clockNowIn(now func() time.Time, loc *time.Location) Clock {
	hr, mm, ss := now().In(loc).Clock()
	return Clock{(hr*60+mm)*60 + ss, loc}
}

// TODO First attempt to parse a timezone
// If no timezone is provided, assume the local timezone
// TODO Replace with a single Must function?

// MustParseClock will panic if the given string cannot be parsed as a clock
// in the local location.
func MustParseClock(value string) Clock {
	clock, err := parseClock(value, time.Local)
	if err != nil {
		panic(err)
	}
	return clock
}

// MustParseClockUTC will panic if the given string cannot be parsed as a
// clock in the UTC location.
func MustParseClockUTC(value string) Clock {
	clock, err := parseClock(value, time.UTC)
	if err != nil {
		panic(err)
	}
	return clock
}

// MustParseClockIn will panic if the given string cannot be parsed as a
// clock in the given location.
func MustParseClockIn(value string, loc *time.Location) Clock {
	clock, err := parseClock(value, loc)
	if err != nil {
		panic(err)
	}
	return clock
}

// ParseClock will attempt to parse the given string as a clock in the local
// location.
func ParseClock(value string) (Clock, error) {
	return parseClock(value, time.Local)
}

// ParseClockUTC will attempt to parse the given string as a clock in the UTC
// location.
func ParseClockUTC(value string) (Clock, error) {
	return parseClock(value, time.UTC)
}

// ParseClockUTC will attempt to parse the given string as a clock in the given
// location.
func ParseClockIn(value string, loc *time.Location) (Clock, error) {
	return parseClock(value, loc)
}

func parseClock(value string, loc *time.Location) (Clock, error) {
	var c Clock
	t, err := time.Parse("15:04:05", value)
	if err != nil {
		return c, err
	}
	hr, mm, ss := t.Clock()
	return Clock{(hr*60+mm)*60 + ss, loc}, nil
}

// Clocks implement the `sort.Interface` for clocks
type Clocks []Clock

// Len returns the length of the clocks slice.
func (c Clocks) Len() int {
	return len(c)
}

// Swap swaps the elements of the clocks slice.
func (c Clocks) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

// Less returns a boolean indicating if the given clocks elements are in
// ascending order.
func (c Clocks) Less(i, j int) bool {
	// Clocks may have been built through addition, so mod by seconds in a day
	return (c[i].sec % (60 * 60 * 24)) < (c[j].sec % (60 * 60 * 24))
}

// SortClocks is a helper method to quickly sort a slice of clocks in ascending
// order.
func SortClocks(clocks []Clock) {
	sort.Sort(Clocks(clocks))
}
