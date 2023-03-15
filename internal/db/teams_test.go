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

func TestTeam_AddPlayer(t *testing.T) {
	Clear()
	Init()
	ctx := context.TODO()
	p1, _ := CreatePlayerWithID(ctx, "12345")
	p2, _ := CreatePlayerWithID(ctx, "12346")
	p3, _ := CreatePlayerWithID(ctx, "12347")
	p4, _ := CreatePlayerWithID(ctx, "12348")
	p5, _ := CreatePlayerWithID(ctx, "12349")
	team1, _ := p1.CreateTeamWithName(ctx, "team1")
	p5.CreateTeamWithName(ctx, "team5")
	t.Run("addoneplayer", func(t *testing.T) {
		err := team1.AddPlayer(ctx, p2)
		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}
	})
	t.Run("addotherteam", func(t *testing.T) {
		err := team1.AddPlayer(ctx, p5)
		if !errors.Is(err, static.ErrUserAlreadyInTeam) {
			t.Errorf("unexpected error, should be: %s", static.ErrUserAlreadyInTeam)
		}
	})
	t.Run("addfull", func(t *testing.T) {
		err := team1.AddPlayer(ctx, p3)
		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}
		err = team1.AddPlayer(ctx, p4)
		if !errors.Is(err, static.ErrTeamFull) {
			t.Errorf("unexpected error, should be: %s", static.ErrTeamFull)
		}
	})
}

func TestTeam_KickPlayer(t *testing.T) {
	Clear()
	Init()
	ctx := context.TODO()
	p1, _ := CreatePlayerWithID(ctx, "12345")
	p2, _ := CreatePlayerWithID(ctx, "12346")
	p3, _ := CreatePlayerWithID(ctx, "12347")
	p4, _ := CreatePlayerWithID(ctx, "12348")
	team1, _ := p1.CreateTeamWithName(ctx, "team1")
	team1.AddPlayer(ctx, p2)
	team1.AddPlayer(ctx, p3)
	t.Run("kickowner", func(t *testing.T) {
		err := team1.KickPlayer(ctx, p1)
		if !errors.Is(err, static.ErrOwnerNotInTeam) {
			t.Errorf("unexpected error, should be: %s", static.ErrOwnerNotInTeam)
		}
	})
	t.Run("kick", func(t *testing.T) {
		err := team1.KickPlayer(ctx, p2)
		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}
	})
	t.Run("kicknobody", func(t *testing.T) {
		err := team1.KickPlayer(ctx, p4)
		if !errors.Is(err, static.ErrPlayerNotInTeam) {
			t.Errorf("unexpected error, should be: %s", static.ErrPlayerNotInTeam)
		}
	})
}

func TestTeam_SetOwner(t *testing.T) {
	Clear()
	Init()
	ctx := context.TODO()
	p1, _ := CreatePlayerWithID(ctx, "12345")
	p2, _ := CreatePlayerWithID(ctx, "12346")
	p3, _ := CreatePlayerWithID(ctx, "12347")
	p4, _ := CreatePlayerWithID(ctx, "12348")
	team1, _ := p1.CreateTeamWithName(ctx, "team1")
	team1.AddPlayer(ctx, p2)
	team1.AddPlayer(ctx, p3)
	t.Run("setownernotinteam", func(t *testing.T) {
		err := team1.SetOwner(ctx, p4)
		if !errors.Is(err, static.ErrOwnerNotInTeam) {
			t.Errorf("unexpected error, should be: %s", static.ErrOwnerNotInTeam)
		}
	})
	t.Run("setowner", func(t *testing.T) {
		err := team1.SetOwner(ctx, p2)
		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}
	})
}

func TestPlayer_LeaveTeam(t *testing.T) {
	Clear()
	Init()
	ctx := context.TODO()
	p, _ := CreatePlayerWithID(ctx, "12345")
	t.Run("noteam", func(t *testing.T) {
		err := p.LeaveTeam(ctx)
		if !errors.Is(err, static.ErrNoTeam) {
			t.Errorf("unexpected error, should be: %s", static.ErrNoTeam)
		}
	})
	p.CreateTeamWithName(ctx, "teamname")
	t.Run("leave", func(t *testing.T) {
		err := p.LeaveTeam(ctx)
		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}
	})
}
