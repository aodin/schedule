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

// Seconds since the beginning of the day
// TODO nano-seconds?
type Clock struct {
	sec int
	loc *time.Location
}

// Pretty printing of the clock
func (c Clock) String() string {
	hour, minute, second := c.HMS()
	return fmt.Sprintf("%d:%02d:%02d", hour, minute, second)
}

func (c Clock) Add(d time.Duration) Clock {
	c.sec += int(d / 1e9)
	return c
}

func (c Clock) Before(other Clock) bool {
	// Clocks may have been built through addition, so mod by seconds in a day
	return (c.sec % (60 * 60 * 24)) < (other.sec % (60 * 60 * 24))
}

func (c Clock) Location() *time.Location {
	return c.loc
}

// TODO Return the clock at a different timezone

func (c Clock) HMS() (int, int, int) {
	return c.Hour(), c.Minute(), c.Second()
}

// If there are more than one day's worth of seconds, roll over to the next day
func (c Clock) Hour() int {
	return (c.sec / (60 * 60)) % 24
}

func (c Clock) Minute() int {
	return (c.sec / 60) % 60
}

func (c Clock) Second() int {
	return c.sec % 60
}

func (c Clock) TotalSeconds() int {
	return c.sec
}

// TODO This operation does not make much sense
func (c Clock) UTC() Clock {
	c.loc = time.UTC
	// TODO Adjust the clock?
	return c
}

// Get the time that represents the next occurence of the clock
// TODO Timezones matter!
func (c Clock) Next() time.Time {
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

func (c Clock) ToTime(y int, m time.Month, d int) time.Time {
	return time.Date(y, m, d, c.Hour(), c.Minute(), c.Second(), 0, c.loc)
}

// Create a Clock from the given time.Time
func ClockFromTime(t time.Time) Clock {
	hr, mm, ss := t.Clock()
	return Clock{(hr*60+mm)*60 + ss, t.Location()}
}

// TODO Interact with the runtime?
func ClockNow() Clock {
	return clockNowIn(defaultNow, time.Local)
}

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
func MustParseClock(value string) Clock {
	clock, err := parseClock(value, time.Local)
	if err != nil {
		panic(err)
	}
	return clock
}

func MustParseClockUTC(value string) Clock {
	clock, err := parseClock(value, time.UTC)
	if err != nil {
		panic(err)
	}
	return clock
}

func MustParseClockIn(value string, loc *time.Location) Clock {
	clock, err := parseClock(value, loc)
	if err != nil {
		panic(err)
	}
	return clock
}

func ParseClock(value string) (Clock, error) {
	return parseClock(value, time.Local)
}

func ParseClockUTC(value string) (Clock, error) {
	return parseClock(value, time.UTC)
}

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

// Implement the sort.Interface for clocks
type Clocks []Clock

func (c Clocks) Len() int {
	return len(c)
}

func (c Clocks) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func (c Clocks) Less(i, j int) bool {
	// Clocks may have been built through addition, so mod by seconds in a day
	return (c[i].sec % (60 * 60 * 24)) < (c[j].sec % (60 * 60 * 24))
}

func SortClocks(clocks []Clock) {
	sort.Sort(Clocks(clocks))
}
