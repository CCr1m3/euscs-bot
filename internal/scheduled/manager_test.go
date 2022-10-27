package scheduled

import (
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

func init() {
	logrus.SetLevel(log.DebugLevel)
}
func Test_scheduledTaskManager_Add(t *testing.T) {
	a := 0
	TaskManager.Add(Task{ID: "increment", Frequency: time.Millisecond * 100, Run: func() {
		a += 1
	}})
	time.Sleep(time.Second)
	TaskManager.Cancel(Task{ID: "increment"})
	if a < 10 {
		t.Errorf("increment was not succesfully ran at least 1000 times: a=%d", a)
	}
	time.Sleep(time.Second)
	if a > 150 {
		t.Errorf("increment was not succesfully stopped: a=%d", a)
	}
}
