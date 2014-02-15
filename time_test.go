package schedule

import (
	"testing"
	"time"
)

// TODO There's no way to manipulate system time?
func TestTime(t *testing.T) {
	eta := time.Now()
	setNow := func() time.Time {
		return eta.Add(-time.Millisecond)
	}
	// TODO Testing channels?
	<-tickAt(setNow, eta)
	// t.Fatal("Tick:", tick, tickAt)
}
