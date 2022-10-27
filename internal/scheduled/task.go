package scheduled

import "time"

type Task struct {
	ID        string
	Run       func()
	Frequency time.Duration
}
