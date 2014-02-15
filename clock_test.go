package schedule

import (
	"testing"
)

// TODO A better place for testing functions?
func expectClock(t *testing.T, clock Clock, hour, min, sec int) {
	if clock.hour != hour {
		t.Errorf("Unexpected clock hours: %d != %d", clock.hour, hour)
	}
	if clock.min != min {
		t.Errorf("Unexpected clock minutes: %d != %d", clock.min, min)
	}
	if clock.sec != sec {
		t.Errorf("Unexpected clock seconds: %d != %d", clock.sec, sec)
	}
}

func expectString(t *testing.T, a, b string) {
	if a != b {
		t.Errorf("Unexpected string: %s != %s", a, b)
	}
}

func TestClock(t *testing.T) {
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
