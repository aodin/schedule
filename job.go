package schedule

import (
	"time"
)

// Job wraps a niladic function that will be performed on every tick for a
// given number of iterations. The easiest way to create a job is through
// the scheduler methods such as Daily and RepeatN.
type Job struct {
	Name      string
	exec      func() error
	quit      chan bool
	tick      <-chan time.Time
	setter    func() <-chan time.Time
	n         int
	increment int
	scheduler *Scheduler
}

// Run will start the job's iteration loop. The job will run on the next tick.
// Jobs are repeated for as many iterations were specified unless the quit
// signal is received. During a job's iteration, the job's parent scheduler
// has its wait group incremented. The wait group is then decremented upon
// completion of the iteration loop or receiving a quit signal.
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

				// Reset the tick if a setter is present
				if j.setter != nil {
					j.tick = j.setter()
				}
			}
		}

		// Remove this job from this scheduler's wait group
		j.scheduler.unfinished.Done()
	}()
}

// Quit will stop the job. If a job is in progress, then it will be completed
// before the job quits.
func (j *Job) Quit() {
	j.quit <- true
}
