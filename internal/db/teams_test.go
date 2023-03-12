package db

import (
	"context"
	"errors"
	"testing"

	"github.com/euscs/euscs-bot/internal/static"
	"github.com/google/go-cmp/cmp"
)

func Test_db_GetTeams(t *testing.T) {
	clearDB()
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
	p1 := &Player{DiscordID: "12345"}
	p1.Save(ctx)
	team := Team{Players: Players{p1}, OwnerID: p1.DiscordID, Name: "teamname"}
	team.Save(ctx)
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
	clearDB()
	Init()
	ctx := context.TODO()
	p1 := &Player{DiscordID: "12345"}
	p1.Save(ctx)
	team := &Team{Players: Players{p1}, OwnerID: p1.DiscordID, Name: "teamname"}
	team.Save(ctx)
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
			t.Logf("want: %#v\n", team)
			t.Logf("got: %#v\n", team2)
			t.Errorf("teams are different")
		}
	})
}

func Test_db_GetTeamByPlayerID(t *testing.T) {
	clearDB()
	Init()
	ctx := context.TODO()
	p1 := &Player{DiscordID: "12345"}
	p1.Save(ctx)
	team := &Team{Players: Players{p1}, OwnerID: p1.DiscordID, Name: "teamname"}
	team.Save(ctx)
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
	clearDB()
	Init()
	ctx := context.TODO()
	p1 := &Player{DiscordID: "12345"}
	p1.Save(ctx)
	team := Team{Players: Players{p1}, OwnerID: p1.DiscordID, Name: "teamname"}

	t.Run("deleteunexistingteam", func(t *testing.T) {
		err := team.Delete(ctx)
		if err != nil {
			t.Errorf("unexpected error while deleting team: %s", err.Error())
		}
	})
	team.Save(ctx)
	t.Run("deleteteam", func(t *testing.T) {
		err := team.Delete(ctx)
		if err != nil {
			t.Errorf("unexpected error while deleting team: %s", err.Error())
		}
	})
}

func TestTeam_Save(t *testing.T) {
	clearDB()
	Init()
	ctx := context.TODO()
	p1 := &Player{DiscordID: "12345"}
	p1.Save(ctx)
	p2 := &Player{DiscordID: "12346"}
	p2.Save(ctx)
	p3 := &Player{DiscordID: "12347"}
	p3.Save(ctx)
	p4 := &Player{DiscordID: "12348"}
	p4.Save(ctx)
	team := Team{Players: Players{p1, p2, p3, p4}, OwnerID: p1.DiscordID, Name: "teamname"}
	t.Run("savewith4players", func(t *testing.T) {
		err := team.Save(ctx)
		if !errors.Is(err, static.ErrTeamFull) {
			t.Errorf("error should be: %s", static.ErrTeamFull)
		}
	})
	team = Team{Players: Players{p1}, OwnerID: p2.DiscordID, Name: "teamname"}
	t.Run("savewithwrongownerid", func(t *testing.T) {
		err := team.Save(ctx)
		if !errors.Is(err, static.ErrOwnerNotInTeam) {
			t.Errorf("error should be: %s", static.ErrOwnerNotInTeam)
		}
	})
	team = Team{Players: Players{p1}, OwnerID: p1.DiscordID, Name: "teamname"}
	t.Run("simplesave", func(t *testing.T) {
		err := team.Save(ctx)
		if err != nil {
			t.Errorf("unexpected error while saving team: %s", err.Error())
		}
	})
	t.Run("saveandedit", func(t *testing.T) {
		err := team.Save(ctx)
		if err != nil {
			t.Errorf("unexpected error while saving team: %s", err.Error())
		}
		team.Players = append(team.Players, p2)
		team.OwnerID = p2.DiscordID
		team.Save(ctx)
		if err != nil {
			t.Errorf("unexpected error while saving team: %s", err.Error())
		}
	})
	team2 := Team{Players: Players{p1}, OwnerID: p1.DiscordID, Name: "teamname2"}
	t.Run("tryingtoaddsomeonealreadyinateam", func(t *testing.T) {
		err := team2.Save(ctx)
		if !errors.Is(err, static.ErrUserAlreadyInTeam) {
			t.Errorf("error should be: %s", static.ErrUserAlreadyInTeam)
		}
	})
	team3 := Team{Players: Players{p3}, OwnerID: p3.DiscordID, Name: "teamname2"}
	team3.Save(ctx)
	t.Run("tryingtoaddsomeonealreadyinateambis", func(t *testing.T) {
		team3.Players = append(team3.Players, p1)
		err := team3.Save(ctx)
		if !errors.Is(err, static.ErrUserAlreadyInTeam) {
			t.Errorf("error should be: %s", static.ErrUserAlreadyInTeam)
		}
	})
	team3 = Team{Players: Players{p3, p4}, OwnerID: p3.DiscordID, Name: "teamname2"}
	team3.Save(ctx)
	t.Run("kickowner", func(t *testing.T) {
		team3.Players = Players{p4}
		err := team3.Save(ctx)
		if !errors.Is(err, static.ErrOwnerNotInTeam) {
			t.Errorf("error should be: %s", static.ErrOwnerNotInTeam)
		}
	})
	t.Run("kickplayer", func(t *testing.T) {
		team3.Players = Players{p3}
		err := team3.Save(ctx)
		if err != nil {
			t.Errorf("unexpected error while saving team: %s", err.Error())
		}
	})
}
