package schedule

import (
	"time"
)

type Job struct {
	Name      string
	exec      func() error
	quit      chan bool
	tick      <-chan time.Time
	setTick   func() <-chan time.Time
	n         int
	increment int
	scheduler *Scheduler
}

func (j *Job) Run() {
	// Add another job to this scheduler's wait group
	j.scheduler.unfinished.Add(1)

	// Perform all iterations of the job in the same goroutine
	go func() {
		// Main iteration loop
	Loop:
		for i := 0; i < j.n; i += j.increment {
			select {
			case <-j.quit:
				// Quit the iteration loop
				break Loop
			case <-j.tick:
				// Run the job and record the time elapsed
				status := Status{Start: time.Now()}
				status.Error = j.exec()
				status.End = time.Now()

				// Send the status to the logger
				j.scheduler.logger.Log(status)

				// Reset the tick
				j.tick = j.setTick()
			}
		}

		// Remove this job from this scheduler's wait group
		j.scheduler.unfinished.Done()
	}()
}

// Allow the job to quit between iterations
func (j *Job) Quit() {
	j.quit <- true
}
