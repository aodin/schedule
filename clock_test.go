package schedule

import (
	"testing"
	"time"
)

// TODO A better place for testing functions?
func expectClock(t *testing.T, clock Clock, hour, min, sec int) {
	// TODO Or just compare the secs and the location
	h, m, s := clock.HMS()
	if h != hour {
		t.Errorf("Unexpected clock hours: %d != %d", h, hour)
	}
	if m != min {
		t.Errorf("Unexpected clock minutes: %d != %d", m, min)
	}
	if s != sec {
		t.Errorf("Unexpected clock seconds: %d != %d", s, sec)
	}
	// TODO Test location as well
}

func expectTime(t *testing.T, a, b time.Time) {
	if a != b {
		t.Errorf("Unexpected time: %s != %s", a, b)
	}
}

func expectString(t *testing.T, a, b string) {
	if a != b {
		t.Errorf("Unexpected string: %s != %s", a, b)
	}
}

func TestClock(t *testing.T) {
	// The the clockNow function with a hardcoded time
	setNow := func() time.Time {
		return time.Date(2014, time.Month(2), 14, 12, 12, 12, 0, time.Local)
	}
	cNow := clockNow(setNow)
	expectClock(t, cNow, 12, 12, 12)

	plusOne := cNow.Add(-time.Minute)
	expectClock(t, plusOne, 12, 11, 12)

	// Get the next occurence of the clock
	fifteenth := plusOne.next(setNow)
	expected := time.Date(2014, time.Month(2), 15, 12, 11, 12, 0, time.Local)
	expectTime(t, fifteenth, expected)

	// Parse some clock strings
	fivePM, err := ParseClock("17:00:00")
	if err != nil {
		t.Fatal(err)
	}
	expectClock(t, fivePM, 17, 0, 0)
	expectString(t, fivePM.String(), "17:00:00")

	fives, err := ParseClock("5:05:05")
	if err != nil {
		t.Fatal(err)
	}
	expectString(t, fives.String(), "5:05:05")
}
