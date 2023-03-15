package scheduled

import "time"

type Task struct {
	ID        string
	Run       func()
	Frequency time.Duration
	At        time.Time
	RunNow    bool
}
