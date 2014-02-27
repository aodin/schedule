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

// Perform a job immediately once
func (s *Scheduler) Now(exec func() error) *Job {
	job := &Job{
		exec:      exec,
		quit:      make(chan bool),
		tick:      time.After(0),
		setTick:   func() <-chan time.Time { return time.After(0) },
		n:         1,
		increment: 1,
		scheduler: s,
	}
	// TODO Some form of job with a delayed run
	job.Run()
	// TODO register the job on the scheduler?
	return job
}

// Perform the job immediately, then repeat forever while waiting the given
// duration between iterations
func (s *Scheduler) Repeat(exec func() error, wait time.Duration) *Job {
	job := &Job{
		exec:      exec,
		quit:      make(chan bool),
		tick:      time.After(0),
		setTick:   func() <-chan time.Time { return time.After(wait) },
		n:         1,
		scheduler: s,
	}
	// TODO Some form of job with a delayed run
	job.Run()
	// TODO register the job on the scheduler
	return job
}

func (s *Scheduler) RepeatN(exec func() error, wait time.Duration, n int) *Job {
	job := &Job{
		exec:      exec,
		quit:      make(chan bool),
		tick:      time.After(0),
		setTick:   func() <-chan time.Time { return time.After(wait) },
		n:         n,
		increment: 1,
		scheduler: s,
	}
	// TODO Some form of job with a delayed run
	job.Run()
	// TODO register the job on the scheduler
	return job
}

func (s *Scheduler) Daily(exec func() error, clock Clock) *Job {
	// Determine the next time the given clock will occur
	next := func() <-chan time.Time { return TickAt(clock.Next()) }
	job := &Job{
		exec:      exec,
		quit:      make(chan bool),
		tick:      next(),
		setTick:   next,
		n:         1,
		scheduler: s,
	}
	// TODO Some form of job with a delayed run
	job.Run()
	// TODO register the job on the scheduler
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

func WaitForJobsToFinish() error {
	return std.WaitForJobsToFinish()
}

// TODO Method to tell all jobs to quit as soon as possible for clean shutdown
