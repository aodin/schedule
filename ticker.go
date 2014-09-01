package schedule

import (
	"sort"
	"time"
)

// Ticker ticks on all the given days / clocks. The easiest way to create a
// ticker is with the functions DayClockTicker and DaysAndClocksTicker.
// TODO daytimes instead of days and clocks?
type Ticker struct {
	C      chan time.Time
	days   []time.Weekday
	clocks []Clock
	loc    *time.Location
	now    func() time.Time
}

func (ticker *Ticker) currentDayIndex(today time.Weekday) int {
	return sort.Search(
		len(ticker.days),
		func(i int) bool { return ticker.days[i] >= today },
	)
}

func (ticker *Ticker) currentClockIndex(clock Clock) int {
	return sort.Search(
		len(ticker.clocks),
		func(i int) bool { return ticker.clocks[i].sec >= clock.sec },
	)
}

// Start the given Ticker.
func (ticker *Ticker) Start() {
	now := ticker.now().In(ticker.loc)

	today := now.Weekday()
	clock := ClockFromTime(now)

	// Find the current position in days and clocks
	d := ticker.currentDayIndex(today)
	c := ticker.currentClockIndex(clock)

	// TODO How to shut down the ticker?
	go func() {
		for {
			// Get the next day that a tick should occur
			// Using mod guarantees the increment will never overflow
			d %= len(ticker.days)
			day := ticker.days[d]

			for c = c; c < len(ticker.clocks); c += 1 {
				nextClock := ticker.clocks[c]

				// Build the next time from this day and clock
				now = ticker.now().In(ticker.loc)
				nowClock := ClockFromTime(now)

				// Difference between days
				// Call the now function as few times as possible
				away := daysAway(func() time.Time { return now }, day)

				// If the next tick day is today, we need to check that
				// the next clock hasn't already occured
				// We perform this check instead next.Before(now), because if
				// we want to delay to tick by a week if the times are equal
				if away == 0 && !nowClock.Before(nextClock) {
					away = 7
				}

				// It is safe to add more days than there are in a month
				// http://golang.org/src/pkg/time/time.go?s=19663:19722#L648
				// TODO Use AddDate
				nextTick := time.Date(
					now.Year(),
					now.Month(),
					now.Day()+away,
					nextClock.Hour(),
					nextClock.Minute(),
					nextClock.Second(),
					0,
					nextClock.Location(),
				)

				tick := <-TickAt(nextTick)
				ticker.C <- tick
			}
			// Reset the clock counter
			c = 0
			d += 1
		}
	}()
}

// Assume the list is sorted and remove any duplicates
func uniqueDays(a []time.Weekday) []time.Weekday {
	if len(a) < 1 {
		return a
	}
	var u time.Weekday
	unique := []time.Weekday{a[0]}
	for i := 1; i < len(a); i += 1 {
		if a[i] != unique[u] {
			unique = append(unique, a[i])
			u += 1
		}
	}
	return unique
}

// DayClockTicker creates a new Ticker that ticks at the time of day specified
// by the clock and for every day in the days array.
func DayClockTicker(weekday time.Weekday, clock Clock) *Ticker {
	// TODO Start the ticker immediately?
	return &Ticker{
		C:      make(chan time.Time),
		days:   []time.Weekday{weekday},
		clocks: []Clock{clock},
		loc:    clock.loc,
		now:    defaultNow,
	}
}

// DaysAndClocksTicker creates a new Ticker that will tick on all
// combinations of the given weekdays and clocks.
func DaysAndClocksTicker(weekdays []time.Weekday, clocks []Clock) *Ticker {
	// Days and clocks may be given out of order, sort them
	SortWeekdays(weekdays)
	SortClocks(clocks)

	// TODO There must be at least one weekday and one clock
	if len(clocks) < 1 {
		return nil
	}

	// TODO Another way to get location?
	loc := clocks[0].loc

	// TODO Start the ticker immediately?
	return &Ticker{
		C:      make(chan time.Time),
		days:   uniqueDays(weekdays), // Remove any duplicate days
		clocks: clocks,               // TODO unique clocks?
		loc:    loc,
		now:    defaultNow,
	}
}
