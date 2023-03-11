package db

import (
	"context"
	"database/sql"

	"github.com/euscs/euscs-bot/internal/models"
)

func CreateTeam(ctx context.Context, t *models.Team) error {
	tx, err := db.Beginx()
	if err != nil {
		return models.ErrDB(err)
	}
	_, err = tx.NamedExec("INSERT INTO teams (name,ownerplayerID) VALUES (:name,:ownerplayerID)", t)
	if err != nil {
		err2 := tx.Rollback()
		if err2 != nil {
			return models.ErrDB(err2)
		}
		return models.ErrDB(err)
	}
	for _, player := range t.Players {
		_, err = tx.Exec("INSERT INTO teamsplayers (team,playerID) VALUES (?,?)", t.Name, player.DiscordID)
		if err != nil {
			err2 := tx.Rollback()
			if err2 != nil {
				return models.ErrDB(err2)
			}
			return models.ErrDB(err)
		}
	}
	err = tx.Commit()
	if err != nil {
		err2 := tx.Rollback()
		if err2 != nil {
			return models.ErrDB(err2)
		}
		return models.ErrDB(err)
	}
	return nil
}

func UpdateTeam(ctx context.Context, t *models.Team) error {
	tx, err := db.Beginx()
	if err != nil {
		return models.ErrDB(err)
	}
	_, err = tx.NamedExec("UPDATE teams set ownerplayerID=:ownerplayerID where name=:name", t)
	if err != nil {
		err2 := tx.Rollback()
		if err2 != nil {
			return models.ErrDB(err2)
		}
		return models.ErrDB(err)
	}
	_, err = tx.NamedExec("DELETE from teamsplayers where team=:name", t)
	if err != nil {
		err2 := tx.Rollback()
		if err2 != nil {
			return models.ErrDB(err2)
		}
		return models.ErrDB(err)
	}
	for _, player := range t.Players {
		_, err = tx.Exec("INSERT INTO teamsplayers (team,playerID) VALUES (?,?)", t.Name, player.DiscordID)
		if err != nil {
			err2 := tx.Rollback()
			if err2 != nil {
				return models.ErrDB(err2)
			}
			return models.ErrDB(err)
		}
	}
	err = tx.Commit()
	if err != nil {
		err2 := tx.Rollback()
		if err2 != nil {
			return models.ErrDB(err2)
		}
		return models.ErrDB(err)
	}
	return nil
}

func GetTeamByName(ctx context.Context, name string) (*models.Team, error) {
	var team models.Team
	err := db.Get(&team, "SELECT name,ownerplayerID FROM teams WHERE name=?", name)
	if err != nil {
		return nil, models.ErrDB(err)
	}
	err = getPlayersInTeam(ctx, &team)
	if err != nil {
		return nil, models.ErrDB(err)
	}
	return &team, nil
}

func GetTeamByPlayerID(ctx context.Context, playerID string) (*models.Team, error) {
	var team models.Team
	err := db.Get(&team, "SELECT name,ownerplayerID FROM teams JOIN teamsplayers ON teamsplayers.team = teams.name WHERE teamsplayers.playerID=?", playerID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, models.ErrNotFound
		}
		return nil, models.ErrDB(err)
	}
	err = getPlayersInTeam(ctx, &team)
	if err != nil {
		return nil, models.ErrDB(err)
	}
	return &team, nil
}

func GetTeams(ctx context.Context) ([]*models.Team, error) {
	teams := []*models.Team{}
	err := db.Select(&teams, "SELECT name,ownerplayerID FROM teams")
	if err != nil {
		return nil, models.ErrDB(err)
	}
	for i := range teams {
		err = getPlayersInTeam(ctx, teams[i])
		if err != nil {
			return nil, models.ErrDB(err)
		}
	}
	return teams, nil
}

func getPlayersInTeam(ctx context.Context, team *models.Team) error {
	players := []*models.Player{}
	err := db.Select(&players, "SELECT elo,discordID,osuser,twitchID FROM players JOIN teamsplayers ON teamsplayers.playerID = players.discordID WHERE team=?", team.Name)
	if err != nil {
		return models.ErrDB(err)
	}
	team.Players = players
	return nil
}
