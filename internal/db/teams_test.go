package db

import (
	"context"
	"errors"
	"testing"

	"github.com/euscs/euscs-bot/internal/static"
	"github.com/google/go-cmp/cmp"
)

func Test_db_GetTeams(t *testing.T) {
	Clear()
	Init()
	ctx := context.TODO()
	t.Run("empty", func(t *testing.T) {
		teams, err := GetTeams(ctx)
		if err != nil {
			t.Errorf("failed to get teams: " + err.Error())
		}
		if len(teams) != 0 {
			t.Errorf("teams is not empty")
		}
	})
	p, _ := CreatePlayerWithID(ctx, "12345")
	p.CreateTeamWithName(ctx, "teamname")
	t.Run("with1team", func(t *testing.T) {
		teams, err := GetTeams(ctx)
		if err != nil {
			t.Errorf("failed to get teams: " + err.Error())
		}
		if len(teams) != 1 {
			t.Errorf("len(teams) should be 1")
		}
	})
}

func Test_db_GetTeamByName(t *testing.T) {
	Clear()
	Init()
	ctx := context.TODO()
	p, _ := CreatePlayerWithID(ctx, "12345")
	team, _ := p.CreateTeamWithName(ctx, "teamname")
	t.Run("wrongname", func(t *testing.T) {
		_, err := GetTeamByName(ctx, "wrongname")
		if !errors.Is(err, static.ErrNotFound) {
			t.Errorf("unexpected error, should be: %s", static.ErrNotFound)
		}
	})
	t.Run("success", func(t *testing.T) {
		team2, err := GetTeamByName(ctx, "teamname")
		if err != nil {
			t.Errorf("unexpected error, %s", err.Error())
		}
		if !cmp.Equal(team, team2) {
			t.Errorf("teams are different: %s", cmp.Diff(team, team2))
		}
	})
}

func Test_db_GetTeamByPlayerID(t *testing.T) {
	Clear()
	Init()
	ctx := context.TODO()
	p, _ := CreatePlayerWithID(ctx, "12345")
	team, _ := p.CreateTeamWithName(ctx, "teamname")
	t.Run("wrongid", func(t *testing.T) {
		_, err := GetTeamByPlayerID(ctx, "123456")
		if !errors.Is(err, static.ErrNotFound) {
			t.Errorf("unexpected error, should be: %s", static.ErrNotFound)
		}
	})
	t.Run("success", func(t *testing.T) {
		team2, err := GetTeamByPlayerID(ctx, "12345")
		if err != nil {
			t.Errorf("unexpected error, %s", err.Error())
		}
		if !cmp.Equal(team, team2) {
			t.Logf("want: %#v\n", team)
			t.Logf("got: %#v\n", team2)
			t.Errorf("teams are different")
		}
	})
}

func TestTeam_Delete(t *testing.T) {
	Clear()
	Init()
	ctx := context.TODO()
	p, _ := CreatePlayerWithID(ctx, "12345")
	team := &Team{Name: "test"}
	t.Run("deleteunexistingteam", func(t *testing.T) {
		err := team.Delete(ctx)
		if !errors.Is(err, static.ErrNotFound) {
			t.Errorf("unexpected error should be: %s", static.ErrNotFound)
		}
	})

	team, _ = p.CreateTeamWithName(ctx, "teamname")
	t.Run("deleteteam", func(t *testing.T) {
		err := team.Delete(ctx)
		if err != nil {
			t.Errorf("unexpected error while deleting team: %s", err.Error())
		}
	})
}

func TestTeam_Save(t *testing.T) {
	Clear()
	Init()
	ctx := context.TODO()
	p1, _ := CreatePlayerWithID(ctx, "12345")
	p2, _ := CreatePlayerWithID(ctx, "12346")
	p3, _ := CreatePlayerWithID(ctx, "12347")
	p4, _ := CreatePlayerWithID(ctx, "12348")
	team1, _ := p1.CreateTeamWithName(ctx, "team1")
	t.Run("savewith4players", func(t *testing.T) {
		team1.Players = Players{p1, p2, p3, p4}
		err := team1.Save(ctx)
		if !errors.Is(err, static.ErrTeamFull) {
			t.Errorf("error should be: %s", static.ErrTeamFull)
		}
	})
	team2, _ := p2.CreateTeamWithName(ctx, "team2")
	t.Run("savewithwrongownerid", func(t *testing.T) {
		team2.OwnerID = p1.DiscordID
		err := team2.Save(ctx)
		if !errors.Is(err, static.ErrOwnerNotInTeam) {
			t.Errorf("error should be: %s", static.ErrOwnerNotInTeam)
		}
	})

	t.Run("saveandedit", func(t *testing.T) {
		team1.Players = Players{p1, p3}
		team1.OwnerID = p3.DiscordID
		err := team1.Save(ctx)
		if err != nil {
			t.Errorf("unexpected error while saving team: %s", err.Error())
		}
	})
	t.Run("tryingtoaddsomeonealreadyinateam", func(t *testing.T) {
		team1.Players = Players{p1, p2, p3}
		err := team1.Save(ctx)
		if !errors.Is(err, static.ErrUserAlreadyInTeam) {
			t.Errorf("error should be: %s", static.ErrUserAlreadyInTeam)
		}
	})
	t.Run("kickowner", func(t *testing.T) {
		team1.Players = Players{p4}
		err := team1.Save(ctx)
		if !errors.Is(err, static.ErrOwnerNotInTeam) {
			t.Errorf("error should be: %s", static.ErrOwnerNotInTeam)
		}
	})
	t.Run("kickplayers", func(t *testing.T) {
		team1.Players = Players{p3}
		err := team1.Save(ctx)
		if err != nil {
			t.Errorf("unexpected error while saving team: %s", err.Error())
		}
	})
}

func TestPlayer_GetTeam(t *testing.T) {
	Clear()
	Init()
	ctx := context.TODO()
	p, _ := CreatePlayerWithID(ctx, "12345")
	t.Run("noteam", func(t *testing.T) {
		_, err := p.GetTeam(ctx)
		if !errors.Is(err, static.ErrNotFound) {
			t.Errorf("unexpected error, should be: %s", static.ErrNotFound)
		}
	})
	team, _ := p.CreateTeamWithName(ctx, "teamname")
	t.Run("success", func(t *testing.T) {
		team2, err := p.GetTeam(ctx)
		if err != nil {
			t.Errorf("unexpected error, %s", err.Error())
		}
		if !cmp.Equal(team, team2) {
			t.Logf("want: %#v\n", team)
			t.Logf("got: %#v\n", team2)
			t.Errorf("teams are different")
		}
	})
}
