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
func Test_scheduledTaskManager_Frequency(t *testing.T) {
	a := 0
	TaskManager.Add(Task{ID: "increment", Frequency: time.Millisecond * 100, Run: func() {
		a += 1
	}})
	time.Sleep(time.Second)
	TaskManager.Cancel(Task{ID: "increment"})
	if a < 10 {
		t.Errorf("increment was not succesfully ran at least 10 times: a=%d", a)
	}
	time.Sleep(time.Second)
	if a > 150 {
		t.Errorf("increment was not succesfully stopped: a=%d", a)
	}
}

func Test_scheduledTaskManager_FrequencyWithStarmTime(t *testing.T) {
	a := 0
	TaskManager.Add(Task{ID: "increment", Frequency: time.Millisecond * 100, At: time.Now().Add(time.Millisecond * 500), Run: func() {
		a += 1
	}})
	time.Sleep(time.Second)
	TaskManager.Cancel(Task{ID: "increment"})
	if a < 5 {
		t.Errorf("increment was not succesfully ran at least 5 times: a=%d", a)
	}
	time.Sleep(time.Second)
	if a > 150 {
		t.Errorf("increment was not succesfully stopped: a=%d", a)
	}
}

func Test_scheduledTaskManager_StartTime(t *testing.T) {
	a := 0
	TaskManager.Add(Task{ID: "increment", At: time.Now().Add(time.Millisecond * 500), Run: func() {
		a += 1
	}})
	time.Sleep(time.Second)
	TaskManager.Cancel(Task{ID: "increment"})
	if a != 1 {
		t.Errorf("increment was not succesfully ran only 1 times: a=%d", a)
	}
	time.Sleep(time.Second)
	if a != 1 {
		t.Errorf("increment was not succesfully stopped: a=%d", a)
	}
}
