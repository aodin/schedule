package schedule

import (
	"sync"
	"time"
)

// Scheduler contains a wait group of all unfinished jobs on the Scheduler and
// an optional Logger.
type Scheduler struct {
	unfinished sync.WaitGroup
	logger     Logger
}

// TODO Options
// * Specify if the job should start immediately
// * Register the job on the scheduler for reporting

// Run the job whenever a tick is received on the time channel
func (s *Scheduler) whenever(exec func() error, tick <-chan time.Time) *Job {
	job := &Job{
		exec:      exec,
		quit:      make(chan bool),
		tick:      tick,
		n:         1,
		scheduler: s,
	}
	return job
}

// Whenever will run the job whenever the job's ticker ticks.
func (s *Scheduler) Whenever(exec func() error, tick <-chan time.Time) *Job {
	job := s.whenever(exec, tick)
	job.Run()
	return job
}

// Every will run the job after every tick of the given duration.
func (s *Scheduler) Every(exec func() error, d time.Duration) *Job {
	return Whenever(exec, time.Tick(d))
}

// Now will run the the job immediately once.
func (s *Scheduler) Now(exec func() error) *Job {
	job := &Job{
		exec:      exec,
		quit:      make(chan bool),
		tick:      time.After(0),
		setter:    func() <-chan time.Time { return time.After(0) },
		n:         1,
		increment: 1,
		scheduler: s,
	}
	job.Run()
	return job
}

// Repeat runs the job immediately, then repeats the job forever while waiting
// the given duration between iterations.
func (s *Scheduler) Repeat(exec func() error, wait time.Duration) *Job {
	job := &Job{
		exec:      exec,
		quit:      make(chan bool),
		tick:      time.After(0),
		setter:    func() <-chan time.Time { return time.After(wait) },
		n:         1,
		scheduler: s,
	}
	job.Run()
	return job
}

// RepeatN runs the job immediately, then repeats the job the given number
// of times, waiting the given duration between iterations.
func (s *Scheduler) RepeatN(exec func() error, wait time.Duration, n int) *Job {
	job := &Job{
		exec:      exec,
		quit:      make(chan bool),
		tick:      time.After(0),
		setter:    func() <-chan time.Time { return time.After(wait) },
		n:         n,
		increment: 1,
		scheduler: s,
	}
	job.Run()
	return job
}

// Daily runs the job once a day at the given clock.
func (s *Scheduler) Daily(exec func() error, clock Clock) *Job {
	// Determine the next time the given clock will occur
	next := func() <-chan time.Time { return TickAt(clock.Next()) }
	job := &Job{
		exec:      exec,
		quit:      make(chan bool),
		tick:      next(),
		setter:    next,
		n:         1,
		scheduler: s,
	}
	job.Run()
	return job
}

// Weekly runs the job on the given weekday and clock.
func (s *Scheduler) Weekly(exec func() error, d time.Weekday, c Clock) *Job {
	// Create the ticker and start immediately
	ticker := DayClockTicker(d, c)
	ticker.Start()

	// Create a job that runs on every tick
	job := s.whenever(exec, ticker.C)
	job.Run()
	return job
}

// DaysAndClocks runs the job on every given combination of the given
// weekdays and clocks.
func (s *Scheduler) DaysAndClocks(exec func() error, ds []time.Weekday, cs []Clock) *Job {
	// Create the ticker
	ticker := DaysAndClocksTicker(ds, cs)
	ticker.Start()

	// Create a job that runs on every tick
	job := s.whenever(exec, ticker.C)
	job.Run()
	return job
}

// WaitForJobsToFinish will wait for all the jobs on the scheduler to finish
// before it returns.
func (s *Scheduler) WaitForJobsToFinish() error {
	s.unfinished.Wait()
	// TODO When would this ever return an error? Allow a timeout to be set?
	return nil
}

// SetLogger allows the Scheduler's Logger to be set.
func (s *Scheduler) SetLogger(l Logger) {
	s.logger = l
}

// New creats a new Scheduler with a default logger.
func New() *Scheduler {
	return &Scheduler{
		logger: &DefaultLogger{},
	}
}

var std = New()

// Whenever will run the job on the default scheduler whenever the job's
// ticker ticks.
func Whenever(exec func() error, tick <-chan time.Time) *Job {
	return std.Whenever(exec, tick)
}

// Every will run the job after every tick of the given duration.
func Every(exec func() error, d time.Duration) *Job {
	return std.Every(exec, d)
}

// Now will run the the job immediately once on the default scheduler.
func Now(exec func() error) *Job {
	return std.Now(exec)
}

// Repeat runs the job immediately on the default scheduler,
// then repeats the job forever while waiting the given duration between
// iterations.
func Repeat(exec func() error, wait time.Duration) *Job {
	return std.Repeat(exec, wait)
}

// RepeatN runs the job immediately on the default scheduler, then repeats the
// job the given number of times, waiting the given duration between
// iterations.
func RepeatN(exec func() error, wait time.Duration, n int) *Job {
	return std.RepeatN(exec, wait, n)
}

// Daily runs the job on the default scheduler once a day at the given clock.
func Daily(exec func() error, clock Clock) *Job {
	return std.Daily(exec, clock)
}

// Weekly runs the job on the default scheduler on the given weekday and clock.
func Weekly(exec func() error, weekday time.Weekday, clock Clock) *Job {
	return std.Weekly(exec, weekday, clock)
}

// DaysAndClocks runs the job on the default scheduler on every given
// combination of the given weekdays and clocks.
func DaysAndClocks(exec func() error, ds []time.Weekday, cs []Clock) *Job {
	return std.DaysAndClocks(exec, ds, cs)
}

// WaitForJobsToFinish will wait for all the jobs on the default scheduler to
// finish before it returns.
func WaitForJobsToFinish() error {
	return std.WaitForJobsToFinish()
}

// TODO Method to tell all jobs to quit as soon as possible for clean shutdown
