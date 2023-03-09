package db

import (
	"context"
	"testing"

	"github.com/euscs/euscs-bot/internal/models"
)

func Test_db_CreateMatch(t *testing.T) {
	clearDB()
	Init()
	ctx := context.TODO()
	err := CreatePlayer(ctx, "12345")
	if err != nil {
		t.Errorf("failed to create player: " + err.Error())
	}
	p1, err := GetPlayerById(ctx, "12345")
	if err != nil {
		t.Errorf("failed to get player: " + err.Error())
	}
	err = CreatePlayer(ctx, "12346")
	if err != nil {
		t.Errorf("failed to create player: " + err.Error())
	}
	p2, err := GetPlayerById(ctx, "12346")
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
	err = CreateMatch(ctx, &match)
	if err != nil {
		t.Errorf("failed to create match: " + err.Error())
	}
	err = CreateMatch(ctx, &match)
	if err == nil {
		t.Errorf("duplicate match created")
	}

	m, err := GetMatchByID(ctx, "id")
	if err != nil || m == nil {
		t.Errorf("failed to get match")
		if err != nil {
			t.Errorf(err.Error())
		}
	}
	m2, err := GetMatchByThreadID(ctx, "threadid")
	if err != nil || m2 == nil {
		t.Errorf("failed to get match")
		if err != nil {
			t.Errorf(err.Error())
		}
	}
	if m.ThreadID != m2.ThreadID {
		t.Errorf("mismatching threadID : %v != %v", m.ThreadID, m2.ThreadID)
	}
	inMatch, err := IsPlayerInMatch(ctx, p1)
	if err != nil || !inMatch {
		t.Errorf("error or player not in match")
		if err != nil {
			t.Errorf(err.Error())
		}
	}
	matches, err := GetRunningMatchesOrderedByTimestamp(ctx)
	if err != nil || len(matches) == 0 {
		t.Errorf("error or no running matches found")
		if err != nil {
			t.Errorf(err.Error())
		}
	}
	m.State = models.MatchStateCanceled
	err = UpdateMatch(ctx, m)
	if err != nil {
		t.Errorf("failed to update match " + err.Error())
	}
	m, err = GetMatchByID(ctx, "id")
	if err != nil || m == nil {
		t.Errorf("failed to get match")
		if err != nil {
			t.Errorf(err.Error())
		}
	}
	if m.State != models.MatchStateCanceled {
		t.Errorf("match did not get canceled")
	}

	inMatch, err = IsPlayerInMatch(ctx, p1)
	if err != nil || inMatch {
		t.Errorf("error or player in match after cancel")
		if err != nil {
			t.Errorf(err.Error())
		}
	}
	matches, err = GetRunningMatchesOrderedByTimestamp(ctx)
	if err != nil || len(matches) != 0 {
		t.Errorf("error or running matches found")
		if err != nil {
			t.Errorf(err.Error())
		}
	}

}
