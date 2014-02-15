package schedule

import (
	"log"
)

type Logger interface {
	Log(Status)
}

type DefaultLogger struct{}

// The DefaultLogger will write the status to log
func (l *DefaultLogger) Log(s Status) {
	log.Println(s)
}
