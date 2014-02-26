schedule
========

Schedule tasks in Go using the `time` and `sync` packages.

The task must be niladic function that returns an error.

```go
func Heartbeat() error {
    return nil
}
```

More complex tasks can be built using methods.

```go
type Task struct {
    count int64
}
func (t *Task) Count() error {
    t.count += 1
    return nil
}
```

To repeat a task every hour for 24 times:

```go
schedule.RepeatN(Heartbeat, time.Hour, 24)
schedule.WaitForJobsToFinish()
```

To repeat a task daily at 3:00 UTC:

```go
threeAM := schedule.MustParseClockUTC("3:00:00")
schedule.Daily(Heartbeat, threeAM)
schedule.WaitForJobsToFinish()
```
