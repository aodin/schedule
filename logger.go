package schedule

import (
	"log"
)

// Logger is an interface for logging status messages.
type Logger interface {
	Log(Status)
}

// DefaultLogger is a simple implementation of the Logger interface that will
// print status messages to the `log` package logger.
type DefaultLogger struct{}

// Log printss status messages to `log` package logger.
func (l *DefaultLogger) Log(s Status) {
	log.Println(s)
}
