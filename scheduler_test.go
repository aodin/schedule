package schedule

import (
	"testing"
	"time"
)

func TestScheduler(t *testing.T) {
	// Variable place holders
	var assertJob *testJob
	var j *Job

	// Run a job immediately one time
	assertJob = newTestJob(t, 1)
	Now(assertJob.Increment)
	WaitForJobsToFinish()
	assertJob.Assert()

	// Run a job immediately 3 times
	assertJob = newTestJob(t, 3)
	RepeatN(assertJob.Increment, time.Millisecond, 3)
	WaitForJobsToFinish()
	assertJob.Assert()

	// Run a job, but quit within a few milliseconds
	// TODO How do we guarantee it ran at least once?
	assertJob = newTestJob(t, 1)
	j = Repeat(assertJob.Increment, time.Second)
	go func() {
		<-time.After(1 * time.Millisecond)
		j.Quit()
	}()
	WaitForJobsToFinish()

	// TODO Test the daily functions
	// These will require some form of clock manipulation, just test that the
	// struct fields were set correctly?
}
