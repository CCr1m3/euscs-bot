package db

import (
	"testing"

	"github.com/haashi/omega-strikers-bot/internal/models"
)

func Test_db_CreateMatch(t *testing.T) {
	clearDB()
	Init()
	err := CreatePlayer("12345")
	if err != nil {
		t.Errorf("failed to create player: " + err.Error())
	}
	p1, err := GetPlayerById("12345")
	if err != nil {
		t.Errorf("failed to get player: " + err.Error())
	}
	err = CreatePlayer("12346")
	if err != nil {
		t.Errorf("failed to create player: " + err.Error())
	}
	p2, err := GetPlayerById("12346")
	if err != nil {
		t.Errorf("failed to get player: " + err.Error())
	}

	match := models.Match{
		Team1:     []*models.Player{p1},
		Team2:     []*models.Player{p2},
		ThreadID:  "threadid",
		MessageID: "messageid",
		ID:        "id",
	}
	err = CreateMatch(&match)
	if err != nil {
		t.Errorf("failed to create match: " + err.Error())
	}
	err = CreateMatch(&match)
	if err == nil {
		t.Errorf("duplicate match created")
	}

	m, err := GetMatchByID("id")
	if err != nil || m == nil {
		t.Errorf("failed to get match")
	}
	m2, err := GetMatchByThreadID("threadid")
	if err != nil || m2 == nil {
		t.Errorf("failed to get match")
	}
	if m.ThreadID != m2.ThreadID {
		t.Errorf("mismatching threadID : %v != %v", m.ThreadID, m2.ThreadID)
	}
	inMatch, err := IsPlayerInMatch(p1)
	if err != nil || !inMatch {
		t.Errorf("error or player not in match")
	}
	matches, err := GetRunningMatchesOrderedByTimestamp()
	if err != nil || len(matches) == 0 {
		t.Errorf("error or no running matches found")
	}
	m.State = models.MatchStateCanceled
	err = UpdateMatch(m)
	if err != nil {
		t.Errorf("failed to update match")
	}
	m, err = GetMatchByID("id")
	if err != nil || m == nil {
		t.Errorf("failed to get match")
	}
	if m.State != models.MatchStateCanceled {
		t.Errorf("match did not get canceled")
	}

	inMatch, err = IsPlayerInMatch(p1)
	if err != nil || inMatch {
		t.Errorf("error or player in match after cancel")
	}
	matches, err = GetRunningMatchesOrderedByTimestamp()
	if err != nil || len(matches) != 0 {
		t.Errorf("error or running matches found")
	}

}
