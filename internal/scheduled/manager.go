package scheduled

import (
	"context"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

var TaskManager scheduledTaskManager

type scheduledTaskManager struct {
	m sync.Map
}

func (g *scheduledTaskManager) Add(t Task) {
	g.Cancel(t)
	ctx, cancel := context.WithCancel(context.Background())
	key := t.ID
	g.m.Store(key, cancel)
	go func() {
		log.Debugf("task %v launched for first launch", t.ID)
		t.Run() //run once
		for {
			select {
			case <-ctx.Done():
				log.Infof("task %v runner was stopped", t.ID)
				return
			case <-time.After(t.Frequency):
				log.Debugf("task %v launched after %s", t.ID, t.Frequency)
				t.Run()
			}
		}
	}()
}

func (g *scheduledTaskManager) Cancel(t Task) {
	key := t.ID
	cancel, exist := g.m.Load(key)
	if exist {
		cancel.(context.CancelFunc)()
	}
}
