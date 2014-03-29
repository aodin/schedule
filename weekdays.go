package schedule

import (
	"sort"
	"time"
)

// Golang's time package specifies the day of the week with Sunday = 0, ...

var Workweek = []time.Weekday{
	time.Monday,
	time.Tuesday,
	time.Wednesday,
	time.Thursday,
	time.Friday,
}

var Weekends = []time.Weekday{time.Sunday, time.Saturday}

// Implement the sort.Interface for weekdays
type Weekdays []time.Weekday

func (w Weekdays) Len() int {
	return len(w)
}

func (w Weekdays) Swap(i, j int) {
	w[i], w[j] = w[j], w[i]
}

func (w Weekdays) Less(i, j int) bool {
	return w[i] < w[j]
}

func SortWeekdays(weekdays []time.Weekday) {
	sort.Sort(Weekdays(weekdays))
}

// How many days are between the current day and the given day of the week?
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
