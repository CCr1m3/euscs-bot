package db

import (
	"context"
	"database/sql"
	"errors"

	"github.com/euscs/euscs-bot/internal/static"
)

type Team struct {
	Players Players
	OwnerID string `db:"ownerplayerID"`
	Name    string `db:"name"`
}

func (t *Team) Delete(ctx context.Context) error {
	_, err := GetTeamByName(ctx, t.Name)
	if err != nil {
		return err
	}
	tx, err := db.Beginx()
	if err != nil {
		return static.ErrDB(err)
	}
	_, err = tx.NamedExec("DELETE from teamsplayers where team=:name", t)
	if err != nil {
		err2 := tx.Rollback()
		if err2 != nil {
			return static.ErrDB(err2)
		}
		return static.ErrDB(err)
	}
	_, err = tx.NamedExec("DELETE FROM teams WHERE name=:name", t)
	if err != nil {
		err2 := tx.Rollback()
		if err2 != nil {
			return static.ErrDB(err2)
		}
		return static.ErrDB(err)
	}
	err = tx.Commit()
	if err != nil {
		err2 := tx.Rollback()
		if err2 != nil {
			return static.ErrDB(err2)
		}
		return static.ErrDB(err)
	}
	return nil
}

func (t *Team) AddPlayer(ctx context.Context, player *Player) error {
	if len(t.Players) >= 3 {
		return static.ErrTeamFull
	}
	team, err := player.GetTeam(ctx)
	if err != nil && !errors.Is(err, static.ErrNotFound) {
		return err
	} else if team != nil {
		return static.ErrUserAlreadyInTeam
	}
	_, err = db.Exec("INSERT INTO teamsplayers (team,playerID) VALUES (?,?)", t.Name, player.DiscordID)
	if err != nil {
		return static.ErrDB(err)
	}
	t.Players = append(t.Players, player)
	return nil
}

func (t *Team) KickPlayer(ctx context.Context, player *Player) error {
	if t.OwnerID == player.DiscordID {
		return static.ErrOwnerNotInTeam
	}
	inTeam := false
	for _, p2 := range t.Players {
		if p2.DiscordID == player.DiscordID {
			inTeam = true
		}
	}
	if !inTeam {
		return static.ErrPlayerNotInTeam
	}
	_, err := db.Exec("DELETE FROM teamsplayers WHERE team=? AND playerID=?", t.Name, player.DiscordID)
	if err != nil {
		return static.ErrDB(err)
	}
	err = getPlayersInTeam(ctx, t)
	if err != nil {
		return err
	}
	return nil
}

func (t *Team) SetOwner(ctx context.Context, player *Player) error {
	inTeam := false
	for _, p2 := range t.Players {
		if p2.DiscordID == player.DiscordID {
			inTeam = true
		}
	}
	if !inTeam {
		return static.ErrOwnerNotInTeam
	}
	_, err := db.Exec("UPDATE teams set ownerplayerID=? where name=?", player.DiscordID, t.Name)
	if err != nil {
		return static.ErrDB(err)
	}
	return nil
}

func (p *Player) CreateTeamWithName(ctx context.Context, teamname string) (*Team, error) {
	_, err := GetTeamByName(ctx, teamname)
	if err != nil && !errors.Is(err, static.ErrNotFound) {
		return nil, err
	} else if err == nil {
		return nil, static.ErrTeamnameTaken
	}
	_, err = GetTeamByPlayerID(ctx, p.DiscordID)
	if err != nil && !errors.Is(err, static.ErrNotFound) {
		return nil, err
	} else if err == nil {
		return nil, static.ErrUserAlreadyInTeam
	}
	tx, err := db.Beginx()
	if err != nil {
		return nil, static.ErrDB(err)
	}
	_, err = tx.Exec("INSERT INTO teams (name,ownerplayerID) VALUES(?,?)", teamname, p.DiscordID)
	if err != nil {
		err2 := tx.Rollback()
		if err2 != nil {
			return nil, static.ErrDB(err2)
		}
		return nil, static.ErrDB(err)
	}
	_, err = tx.Exec("INSERT INTO teamsplayers (team,playerID) VALUES (?,?)", teamname, p.DiscordID)
	if err != nil {
		err2 := tx.Rollback()
		if err2 != nil {
			return nil, static.ErrDB(err2)
		}
		return nil, static.ErrDB(err)
	}
	err = tx.Commit()
	if err != nil {
		err2 := tx.Rollback()
		if err2 != nil {
			return nil, static.ErrDB(err2)
		}
		return nil, static.ErrDB(err)
	}
	return &Team{OwnerID: p.DiscordID, Players: Players{p}, Name: teamname}, nil
}

func GetTeamByName(ctx context.Context, name string) (*Team, error) {
	var team Team
	err := db.Get(&team, "SELECT name,ownerplayerID FROM teams WHERE name=?", name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, static.ErrNotFound
		}
		return nil, static.ErrDB(err)
	}
	err = getPlayersInTeam(ctx, &team)
	if err != nil {
		return nil, err
	}
	return &team, nil
}

func GetTeamByPlayerID(ctx context.Context, playerID string) (*Team, error) {
	var team Team
	err := db.Get(&team, "SELECT name,ownerplayerID FROM teams JOIN teamsplayers ON teamsplayers.team = teams.name WHERE teamsplayers.playerID=?", playerID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, static.ErrNotFound
		}
		return nil, static.ErrDB(err)
	}
	err = getPlayersInTeam(ctx, &team)
	if err != nil {
		return nil, err
	}
	return &team, nil
}

func (p *Player) GetTeam(ctx context.Context) (*Team, error) {
	var team Team
	err := db.Get(&team, "SELECT name,ownerplayerID FROM teams JOIN teamsplayers ON teamsplayers.team = teams.name WHERE teamsplayers.playerID=?", p.DiscordID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, static.ErrNotFound
		}
		return nil, static.ErrDB(err)
	}
	err = getPlayersInTeam(ctx, &team)
	if err != nil {
		return nil, err
	}
	return &team, nil
}

func (p *Player) LeaveTeam(ctx context.Context) error {
	team, err := p.GetTeam(ctx)
	if err != nil {
		if errors.Is(err, static.ErrNotFound) {
			return static.ErrNoTeam
		}
		return err
	}
	if p.DiscordID == team.OwnerID {
		return team.Delete(ctx)
	} else {
		return team.KickPlayer(ctx, p)
	}
}

func GetTeams(ctx context.Context) ([]*Team, error) {
	teams := []*Team{}
	err := db.Select(&teams, "SELECT name,ownerplayerID FROM teams")
	if err != nil {
		return nil, static.ErrDB(err)
	}
	for i := range teams {
		err = getPlayersInTeam(ctx, teams[i])
		if err != nil {
			return nil, err
		}
	}
	return teams, nil
}

func getPlayersInTeam(ctx context.Context, team *Team) error {
	players := Players{}
	err := db.Select(&players, "SELECT elo,discordID,osuser,twitchID FROM players JOIN teamsplayers ON teamsplayers.playerID = players.discordID WHERE team=?", team.Name)
	if err != nil {
		return static.ErrDB(err)
	}
	team.Players = players
	return nil
}
