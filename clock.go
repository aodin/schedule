package schedule

import (
	"fmt"
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

// This operation does not make much sense
func (c Clock) UTC() Clock {
	c.loc = time.UTC
	// TODO Adjust the clock
	return c
}

// Get the time that represents the next occurence of the clock
// TODO Timezones matter!
func (c Clock) Next() time.Time {
	return c.next(now)
}

func (c Clock) next(getNow func() time.Time) time.Time {
	n := getNow()
	year, month, day := n.Date()
	hr, mm, ss := c.HMS()
	nxt := time.Date(year, month, day, hr, mm, ss, 0, c.Location())
	if nxt.Before(n) {
		// If nxt has already occured, add a day
		nxt = nxt.Add(24 * time.Hour)
	}
	return nxt
}

func (c Clock) ToTime(y int, m time.Month, d int) time.Time {
	return time.Date(y, m, d, c.Hour(), c.Minute(), c.Second(), 0, c.loc)
}

// TODO Interact with the runtime?
func ClockNow() Clock {
	return clockNow(now)
}

func clockNow(getNow func() time.Time) Clock {
	hr, mm, ss := getNow().Clock()
	return Clock{(hr*60+mm)*60 + ss, time.Local}
}

func ClockNowUTC() Clock {
	return clockNowUTC(now)
}

func clockNowUTC(getNow func() time.Time) Clock {
	hr, mm, ss := getNow().UTC().Clock()
	return Clock{(hr*60+mm)*60 + ss, time.Local}
}

// TODO First attempt to parse a timezone
// If no timezone is provided, assume the local timezone
func ParseClock(value string) (Clock, error) {
	return parseClock(value, time.Local)
}

func ParseClockUTC(value string) (Clock, error) {
	return parseClock(value, time.UTC)
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
