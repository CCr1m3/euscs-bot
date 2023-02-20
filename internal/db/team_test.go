package db

import (
	"context"
	"testing"

	"github.com/haashi/omega-strikers-bot/internal/models"
)

func Test_db_CreateTeam(t *testing.T) {
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
	err = CreatePlayer(ctx, "12347")
	if err != nil {
		t.Errorf("failed to create player: " + err.Error())
	}
	p3, err := GetPlayerById(ctx, "12347")
	if err != nil {
		t.Errorf("failed to get player: " + err.Error())
	}

	team := &models.Team{
		Players: []*models.Player{p1, p2},
		Name:    "teamname",
		OwnerID: p1.DiscordID,
	}

	err = CreateTeam(ctx, team)
	if err != nil {
		t.Errorf("failed to create team: " + err.Error())
	}
	err = CreateTeam(ctx, team)
	if err == nil {
		t.Errorf("duplicate team created")
	}
	team, err = GetTeamByName(ctx, "teamname")
	if err != nil {
		t.Errorf("failed to get team: " + err.Error())
	}
	if len(team.Players) != 2 {
		t.Errorf("failed to create team with 2 members")
	}
	team.Players = []*models.Player{p1}
	err = UpdateTeam(ctx, team)
	if err != nil {
		t.Errorf("failed to update team: " + err.Error())
	}
	team, err = GetTeamByName(ctx, "teamname")
	if err != nil {
		t.Errorf("failed to get team: " + err.Error())
	}
	if len(team.Players) != 1 {
		t.Errorf("failed to remove a member of a team")
	}

	team2 := &models.Team{
		Players: []*models.Player{p1},
		Name:    "team2",
		OwnerID: p1.DiscordID,
	}
	err = CreateTeam(ctx, team2)
	if err == nil {
		t.Errorf("able to create team with same owner")
	}
	team2 = &models.Team{
		Players: []*models.Player{p3},
		Name:    "team2",
		OwnerID: p3.DiscordID,
	}
	err = CreateTeam(ctx, team2)
	if err != nil {
		t.Errorf("failed to create team: " + err.Error())
	}
	team2, err = GetTeamByName(ctx, "team2")
	if err != nil {
		t.Errorf("failed to get team: " + err.Error())
	}
	team2.Players = append(team2.Players, p1)
	err = UpdateTeam(ctx, team2)
	if err == nil {
		t.Errorf("able to add player with team on a new team")
	}
}
