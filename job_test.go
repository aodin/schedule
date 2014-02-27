package schedule

import (
	"testing"
)

// A job that can assert the number of times it has been run
type testJob struct {
	t     *testing.T
	count int
	expected int
}

func (t *testJob) Assert() {
	if t.expected != t.count {
		t.t.Errorf("The job was expected to run %d times, but ran %d", t.expected, t.count)
	}
}

// The job that will be run
func (t *testJob) Increment() error {
	t.count += 1
	return nil
}

func newTestJob(t *testing.T, n int) *testJob {
	return &testJob{t, 0, n} 
}
