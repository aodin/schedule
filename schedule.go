package schedule

import (
	"sync"
	"time"
)

// Maintain a wait group of all unfinished jobs on the queue
type Queue struct {
	unfinished sync.WaitGroup
	logger     Logger
}

// TODO Replace job with an interface?
// TODO Create a common Task struct that can handle these cases?
func (q *Queue) Now(job func() error) {
	// Do the job immediately
	q.unfinished.Add(1)
	go func() {
		// Create a status and record the start time
		// TODO How to record the original job?
		s := Status{
			Start: time.Now(),
		}

		// TODO Should tasks all implement a Job interface
		// In that case, this operation must return an error
		s.Error = job()
		s.End = time.Now()
		if s.Error == nil {
			s.OK = true
		}

		// Send the status to the logger
		q.logger.Log(s)

		// Decrement the job counter
		q.unfinished.Done()
	}()
}

func (q *Queue) Repeat(job func() error, freq time.Duration) {
	q.unfinished.Add(1)
	// TODO Create a new go func for each iteration of the job
	// if the job needs to be STARTED every duration
	go func() {
		for {
			s := Status{
				Start: time.Now(),
			}
			// TODO Should tasks all implement a Job interface
			// In that case, this operation must return an error
			s.Error = job()
			s.End = time.Now()
			if s.Error == nil {
				s.OK = true
			}
			// Send the status to the logger
			q.logger.Log(s)
			time.Sleep(freq)
		}
		q.unfinished.Done()
	}()
}

func (q *Queue) RepeatN(job func() error, freq time.Duration, n int) {
	q.unfinished.Add(1)
	// TODO Create a new go func for each iteration of the job
	// if the job needs to be STARTED every duration
	go func() {
		for i := 0; i < n; i++ {
			s := Status{
				Start: time.Now(),
			}
			// TODO Should tasks all implement a Job interface
			// In that case, this operation must return an error
			s.Error = job()
			s.End = time.Now()
			if s.Error == nil {
				s.OK = true
			}
			// Send the status to the logger
			q.logger.Log(s)
			time.Sleep(freq)
		}
		q.unfinished.Done()
	}()
}

// Perform the task daily at the given time of day
func (q *Queue) Daily(job func() error, clock Clock) {
	q.unfinished.Add(1)
	// TODO New goroutine for each iteration?
	go func() {
		for {
			// Wait for this next time to occur
			<-TickAt(clock.Next())
			s := Status{
				Start: time.Now(),
			}
			// TODO Should tasks all implement a Job interface
			// In that case, this operation must return an error
			s.Error = job()
			s.End = time.Now()
			if s.Error == nil {
				s.OK = true
			}
			// Send the status to the logger
			q.logger.Log(s)
		}
		q.unfinished.Done()
	}()
}

// TODO When would this ever return an error?
func (q *Queue) WaitForJobsToFinish() error {
	q.unfinished.Wait()
	return nil
}

// Create a queue with a default logger
func New() *Queue {
	return &Queue{
		logger: &DefaultLogger{},
	}
}

var std = New()

// Aliases for operations on the std queue
func Now(job func() error) {
	std.Now(job)
}

func Repeat(job func() error, freq time.Duration) {
	std.Repeat(job, freq)
}

func RepeatN(job func() error, freq time.Duration, n int) {
	std.RepeatN(job, freq, n)
}

func Daily(job func() error, clock Clock) {
	std.Daily(job, clock)
}

func WaitForJobsToFinish() error {
	return std.WaitForJobsToFinish()
}
