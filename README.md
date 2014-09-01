schedule
========

Schedule jobs in Go using the `time` and `sync` packages.

Jobs must be niladic functions that return an error.

```go
func Heartbeat() error {
    return nil
}
```

More complex jobs can be built using methods.

```go
type Counter struct {
    count int64
}
func (c *Counter) Count() error {
    c.count += 1
    return nil
}
```

To repeat a job every hour for 24 times:

```go
schedule.RepeatN(Heartbeat, time.Hour, 24)
schedule.WaitForJobsToFinish()
```

To repeat a job daily at 3:00 UTC:

```go
threeAM := schedule.MustParseClockUTC("3:00:00")
schedule.Daily(Heartbeat, threeAM)
schedule.WaitForJobsToFinish()
```

To stop a job cleanly between iterations while it is running forever:

```go
job := schedule.Repeat(c.Count, time.Second)
go func() {
    <- time.After(5 * time.Second)
    job.Quit()
}()
schedule.WaitForJobsToFinish()
```

For full API documentation visit the project's [GoDoc page](https://godoc.org/github.com/aodin/schedule).

-aodin, 2014