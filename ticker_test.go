package schedule

import (
	"testing"
	"time"
)

func expectWeekdayArray(t *testing.T, a, b []time.Weekday) {
	if len(a) != len(b) {
		t.Fatalf("Unexpected array length: %d != %d", len(a), len(b))
	}
	for i, elem := range a {
		if b[i] != elem {
			t.Errorf("Unexpected element at index %d: %d != %d", i, elem, b[i])
		}
	}
}

func TestUniqueDays(t *testing.T) {
	raw := []time.Weekday{1, 2, 3, 3, 3, 4, 5, 5}
	expected := []time.Weekday{1, 2, 3, 4, 5}
	expectWeekdayArray(t, uniqueDays(raw), expected)
}

// TODO Naming scheme for struct methods?
func TestCurrentTickerIndexes(t *testing.T) {
	ticker := &Ticker{
		days: Workweek,
	}
	d := ticker.currentDayIndex(time.Tuesday)
	if d != 1 {
		t.Errorf("Unexpect day index: %d != 1", d)
	}
}
