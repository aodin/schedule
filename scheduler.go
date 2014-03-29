package schedule

import (
	"sync"
	"time"
)

// Maintain a wait group of all unfinished jobs on the Scheduler
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

func (s *Scheduler) Whenever(exec func() error, tick <-chan time.Time) *Job {
	job := s.whenever(exec, tick)
	job.Run()
	return job
}

// Run the job after every tick of the given duration
func (s *Scheduler) Every(exec func() error, d time.Duration) *Job {
	return Whenever(exec, time.Tick(d))
}

// Perform a job immediately once
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

// Perform the job immediately, then repeat forever while waiting the given
// duration between iterations
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

// Run a job on the scheduler that will run on the given weekday and clock
func (s *Scheduler) Weekly(exec func() error, d time.Weekday, c Clock) *Job {
	// Create the ticker and start immediately
	ticker := DayClockTicker(d, c)
	ticker.Start()

	// Create a job that runs on every tick
	job := s.whenever(exec, ticker.C)
	job.Run()
	return job
}

func (s *Scheduler) DaysAndClocks(exec func() error, ds []time.Weekday, cs []Clock) *Job {
	// Create the ticker
	ticker := DaysAndClocksTicker(ds, cs)
	ticker.Start()

	// Create a job that runs on every tick
	job := s.whenever(exec, ticker.C)
	job.Run()
	return job
}

// TODO When would this ever return an error? Allow a timeout to be set?
func (s *Scheduler) WaitForJobsToFinish() error {
	s.unfinished.Wait()
	return nil
}

func (s *Scheduler) SetLogger(l Logger) {
	s.logger = l
}

// Create a Scheduler with a default logger
func New() *Scheduler {
	return &Scheduler{
		logger: &DefaultLogger{},
	}
}

var std = New()

// Aliases for operations on the std Scheduler
func Whenever(exec func() error, tick <-chan time.Time) *Job {
	return std.Whenever(exec, tick)
}

func Every(exec func() error, d time.Duration) *Job {
	return std.Every(exec, d)
}

func Now(exec func() error) *Job {
	return std.Now(exec)
}

func Repeat(exec func() error, wait time.Duration) *Job {
	return std.Repeat(exec, wait)
}

func RepeatN(exec func() error, wait time.Duration, n int) *Job {
	return std.RepeatN(exec, wait, n)
}

func Daily(exec func() error, clock Clock) *Job {
	return std.Daily(exec, clock)
}

func Weekly(exec func() error, weekday time.Weekday, clock Clock) *Job {
	return std.Weekly(exec, weekday, clock)
}

func DaysAndClocks(exec func() error, ds []time.Weekday, cs []Clock) *Job {
	return std.DaysAndClocks(exec, ds, cs)
}

func WaitForJobsToFinish() error {
	return std.WaitForJobsToFinish()
}

// TODO Method to tell all jobs to quit as soon as possible for clean shutdown
