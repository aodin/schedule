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

// TODO Seconds since the beginning of the day
// TODO Location or timezone?
// TODO nano-seconds?
type Clock struct {
	hour int
	min  int
	sec  int
}

func (c Clock) String() string {
	return fmt.Sprintf("%d:%02d:%02d", c.hour, c.min, c.sec)
}

func ParseClock(value string) (Clock, error) {
	var c Clock
	t, err := time.Parse("15:04:05", value)
	if err != nil {
		return c, err
	}
	c.hour, c.min, c.sec = t.Clock()
	return c, nil
}
