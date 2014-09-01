package schedule

import (
	"sort"
	"time"
)

// This file defines convenience types and functions for operating on days.
// Golang's time package specifies the day of the week with Sunday = 0, ...

// Workweek is a slice of the traditional workweek: Monday thru Friday.
var Workweek = []time.Weekday{
	time.Monday,
	time.Tuesday,
	time.Wednesday,
	time.Thursday,
	time.Friday,
}

// Weekends is a slice of the traditional weekend: Saturday and Sunday
var Weekends = []time.Weekday{time.Sunday, time.Saturday}

// Weekdays implements the sort.Interface for Weekdays
type Weekdays []time.Weekday

// Len returns length of the Weekdays slice
func (w Weekdays) Len() int {
	return len(w)
}

// Swap swaps elements in a Weekdays slice
func (w Weekdays) Swap(i, j int) {
	w[i], w[j] = w[j], w[i]
}

// Less determines ordering in a Weekdays slice
func (w Weekdays) Less(i, j int) bool {
	return w[i] < w[j]
}

// SortWeekdays sorts a slice of Weekdays in ascending order (Monday first)
func SortWeekdays(weekdays []time.Weekday) {
	sort.Sort(Weekdays(weekdays))
}

// Days away returns the number of days between the current day and the
// next occurrence if the given Weekday.
func DaysAway(weekday time.Weekday) int {
	return daysAway(defaultNow, weekday)
}

func daysAway(now func() time.Time, weekday time.Weekday) int {
	days := (int(weekday) - int(now().Weekday())) % 7
	if days < 0 {
		days = 7 + days
	}
	return days
}
