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

func (t *Team) Save(ctx context.Context) error {
	if len(t.Players) > 3 {
		return static.ErrTeamFull
	}
	ownerInTeam := false
	for _, player := range t.Players {
		if t.OwnerID == player.DiscordID {
			ownerInTeam = true
		}
	}
	if !ownerInTeam {
		return static.ErrOwnerNotInTeam
	}
	currentTeam, err := GetTeamByName(ctx, t.Name)
	if err != nil && !errors.Is(err, static.ErrNotFound) {
		return err
	} else if errors.Is(err, static.ErrNotFound) {
		tx, err := db.Beginx()
		if err != nil {
			return static.ErrDB(err)
		}
		_, err = tx.NamedExec("INSERT INTO teams (name,ownerplayerID) VALUES(:name,:ownerplayerID)", t)
		if err != nil {
			err2 := tx.Rollback()
			if err2 != nil {
				return static.ErrDB(err2)
			}
			return static.ErrDB(err)
		}
		for _, player := range t.Players {
			team, err := GetTeamByPlayerID(ctx, player.DiscordID)
			if err != nil && !errors.Is(err, static.ErrNotFound) {
				err2 := tx.Rollback()
				if err2 != nil {
					return static.ErrDB(err2)
				}
				return static.ErrDB(err)
			} else if team != nil {
				err2 := tx.Rollback()
				if err2 != nil {
					return static.ErrDB(err2)
				}
				return static.ErrUserAlreadyInTeam
			} else {
				_, err = tx.Exec("INSERT INTO teamsplayers (team,playerID) VALUES(?,?)", t.Name, player.DiscordID)
				if err != nil {
					err2 := tx.Rollback()
					if err2 != nil {
						return static.ErrDB(err2)
					}
					return static.ErrDB(err)
				}
			}
		}
		err = tx.Commit()
		if err != nil {
			err2 := tx.Rollback()
			if err2 != nil {
				return static.ErrDB(err2)
			}
			return static.ErrDB(err)
		}
	} else {
		tx, err := db.Beginx()
		if err != nil {
			return static.ErrDB(err)
		}
		if currentTeam.OwnerID != t.OwnerID {
			_, err = tx.NamedExec("UPDATE teams set ownerplayerID=:ownerplayerID where name=:name", t)
			if err != nil {
				err2 := tx.Rollback()
				if err2 != nil {
					return static.ErrDB(err2)
				}
				return static.ErrDB(err)
			}
		}
		for _, currentPlayer := range currentTeam.Players {
			playerKicked := true
			for _, player := range t.Players {
				if currentPlayer.DiscordID == player.DiscordID {
					playerKicked = false
				}
			}
			if playerKicked {
				_, err = tx.Exec("DELETE FROM teamsplayers WHERE team=? AND playerID=?", t.Name, currentPlayer.DiscordID)
				if err != nil {
					err2 := tx.Rollback()
					if err2 != nil {
						return static.ErrDB(err2)
					}
					return static.ErrDB(err)
				}
			}
		}

		for _, player := range t.Players {
			newPlayer := true
			for _, currentPlayer := range currentTeam.Players {
				if currentPlayer.DiscordID == player.DiscordID {
					newPlayer = false
				}
			}
			if newPlayer {
				team, err := GetTeamByPlayerID(ctx, player.DiscordID)
				if err != nil && !errors.Is(err, static.ErrNotFound) {
					err2 := tx.Rollback()
					if err2 != nil {
						return static.ErrDB(err2)
					}
					return static.ErrDB(err)
				} else if team != nil {
					err2 := tx.Rollback()
					if err2 != nil {
						return static.ErrDB(err2)
					}
					return static.ErrUserAlreadyInTeam
				} else {
					_, err = tx.Exec("INSERT INTO teamsplayers (team,playerID) VALUES (?,?)", t.Name, player.DiscordID)
					if err != nil {
						err2 := tx.Rollback()
						if err2 != nil {
							return static.ErrDB(err2)
						}
						return static.ErrDB(err)
					}
				}
			}
		}
		err = tx.Commit()
		if err != nil {
			err2 := tx.Rollback()
			if err2 != nil {
				return static.ErrDB(err2)
			}
			return static.ErrDB(err)
		}
	}
	return nil
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
		return nil, static.ErrDB(err)
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
		return nil, static.ErrDB(err)
	}
	return &team, nil
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
			return nil, static.ErrDB(err)
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
