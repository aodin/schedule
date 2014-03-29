package schedule

import (
	"testing"
	"time"
)

// TODO A better place for testing functions?
func expectClock(t *testing.T, clock Clock, hour, min, sec int) {
	// TODO Or just compare the secs
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
}

func expectLocation(t *testing.T, a, b *time.Location) {
	if a != b {
		t.Errorf("Unexpected location: %s != %s", a, b)
	}
}

func expectTime(t *testing.T, a, b time.Time) {
	if a != b {
		t.Errorf("Unexpected time: %s != %s", a, b)
	}
}

func expectInt(t *testing.T, a, b int) {
	if a != b {
		t.Errorf("Unexpected integer: %d != %d", a, b)
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

	// Test that the UTC functions return UTC
	cNowUTC := ClockNowUTC()
	expectLocation(t, cNowUTC.loc, time.UTC)

	cNow := ClockNow()
	expectLocation(t, cNow.loc, time.Local)

	// Test Local
	clock := clockNowIn(setNow, time.Local)
	expectClock(t, clock, 12, 12, 12)
	expectLocation(t, clock.loc, time.Local)

	plusOne := clock.Add(-time.Minute)
	expectClock(t, plusOne, 12, 11, 12)
	expectLocation(t, plusOne.loc, time.Local)

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

func TestClock_Before(t *testing.T) {
	threePM := MustParseClock("15:00:00")
	sixPM := MustParseClock("18:00:00")
	if !threePM.Before(sixPM) {
		t.Error("Three PM should be before six PM")
	}
	if threePM.Before(threePM) {
		t.Error("Three PM should not be before itself")
	}
	if sixPM.Before(threePM) {
		t.Error("Six PM should not be before three PM")
	}

	// Build a clock through addition
	twoAM := sixPM.Add(8 * time.Hour)
	if !twoAM.Before(sixPM) {
		t.Error("Two AM should be before six PM")
	}
}
