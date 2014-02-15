package schedule

import ()

type Runnable interface {
	Run() error
}
