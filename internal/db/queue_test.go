package db

import (
	"testing"

	"github.com/haashi/omega-strikers-bot/internal/models"
)

func Test_db_AddPlayerToQueue(t *testing.T) {
	clearDB()
	Init()
	err := CreatePlayer("12345")
	if err != nil {
		t.Errorf("failed to create player: " + err.Error())
	}
	p, err := GetPlayerById("12345")
	if err != nil {
		t.Errorf("failed to get player: " + err.Error())
	}

	if err := AddPlayerToQueue(p, models.RoleFlex); err != nil {
		t.Errorf("AddPlayerToQueue() error: " + err.Error())
	}
	if inQueue, err := IsPlayerInQueue(p); err != nil || !inQueue {
		t.Errorf("AddPlayerToQueue() error: player is not in queue" + err.Error())
	}
	if err := AddPlayerToQueue(p, models.RoleFlex); err == nil {
		t.Errorf("AddPlayerToQueue() should be in error: player is already in queue")
	}
	ps, err := GetPlayersInQueue()
	if err != nil {
		t.Errorf("failed to fetch players in queue: " + err.Error())
	}
	if len(ps) == 0 {
		t.Errorf("failed to fetch players in queue: no players in queue")
	}
	if ps[0].Role != models.RoleFlex {
		t.Errorf("players was queued with the wrong role: " + string(ps[0].Role))
	}
}

func Test_db_RemovePlayerFromQueue(t *testing.T) {
	clearDB()
	Init()
	err := CreatePlayer("12345")
	if err != nil {
		t.Errorf("failed to create player: " + err.Error())
	}
	p, err := GetPlayerById("12345")
	if err != nil {
		t.Errorf("failed to get player: " + err.Error())
	}

	if err := AddPlayerToQueue(p, models.RoleFlex); err != nil {
		t.Errorf("AddPlayerToQueue() error: " + err.Error())
	}
	if err := RemovePlayerFromQueue(p); err != nil {
		t.Errorf("RemovePlayerFromQueue() error: ")
	}
	if inQueue, err := IsPlayerInQueue(p); err != nil || inQueue {
		t.Errorf("RemovePlayerFromQueue() error: player is in queue" + err.Error())
	}

}

func Test_db_GetGoaliesCountInQueue(t *testing.T) {
	clearDB()
	Init()
	got, err := GetGoaliesCountInQueue()
	if err != nil {
		t.Errorf("GetGoaliesCountInQueue() error: " + err.Error())
		return
	}
	if got != 0 {
		t.Errorf("GetGoaliesCountInQueue() = %v, but nobody is in queue yet", got)
	}
	err = CreatePlayer("12345")
	if err != nil {
		t.Errorf("failed to create player: " + err.Error())
	}
	p, err := GetPlayerById("12345")
	if err != nil {
		t.Errorf("failed to get player: " + err.Error())
	}

	if err := AddPlayerToQueue(p, models.RoleFlex); err != nil {
		t.Errorf("AddPlayerToQueue() error: " + err.Error())
	}
	err = CreatePlayer("12346")
	if err != nil {
		t.Errorf("failed to create player: " + err.Error())
	}
	p2, err := GetPlayerById("12346")
	if err != nil {
		t.Errorf("failed to get player: " + err.Error())
	}

	if err := AddPlayerToQueue(p2, models.RoleGoalie); err != nil {
		t.Errorf("AddPlayerToQueue() error: " + err.Error())
	}

	got, err = GetGoaliesCountInQueue()
	if err != nil {
		t.Errorf("GetGoaliesCountInQueue() error: " + err.Error())
		return
	}
	if got != 2 {
		t.Errorf("GetGoaliesCountInQueue() != %v, but one flex and one goalie in queue", got)
	}
}

func Test_db_GetForwardsCountInQueue(t *testing.T) {
	clearDB()
	Init()
	got, err := GetForwardsCountInQueue()
	if err != nil {
		t.Errorf("GetForwardsCountInQueue() error: " + err.Error())
		return
	}
	if got != 0 {
		t.Errorf("GetForwardsCountInQueue() = %v, but nobody is in queue yet", got)
	}
	err = CreatePlayer("12345")
	if err != nil {
		t.Errorf("failed to create player: " + err.Error())
	}
	p, err := GetPlayerById("12345")
	if err != nil {
		t.Errorf("failed to get player: " + err.Error())
	}

	if err := AddPlayerToQueue(p, models.RoleForward); err != nil {
		t.Errorf("AddPlayerToQueue() error: " + err.Error())
	}

	got, err = GetForwardsCountInQueue()
	if err != nil {
		t.Errorf("GetForwardsCountInQueue() error: " + err.Error())
		return
	}
	if got != 1 {
		t.Errorf("GetForwardsCountInQueue() != %v, but one flex and one goalie in queue", got)
	}
}
