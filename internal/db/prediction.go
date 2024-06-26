package db

import (
	"context"

	"github.com/euscs/euscs-bot/internal/static"
	log "github.com/sirupsen/logrus"
)

type Prediction struct {
	Player
	Team   int `db:"team"`
	Amount int `db:"amount"`
}

func GetPlayersPredictionOnMatch(ctx context.Context, match *Match) ([]*Prediction, error) {
	predictions := []*Prediction{}
	err := db.Select(&predictions, "SELECT elo,discordID,osuser,lastrankupdate,team,credits,amount FROM players JOIN predictions ON predictions.playerID = players.discordID WHERE matchID=? ", match.ID)
	if err != nil {
		return nil, static.ErrDB(err)
	}
	return predictions, nil
}

func CreatePrediction(ctx context.Context, discordID string, matchID string, team int, amount int) error {
	tx, err := db.Beginx()
	if err != nil {
		log.WithFields(log.Fields{
			string(static.UUIDKey):  ctx.Value(static.UUIDKey),
			string(static.ErrorKey): err.Error(),
		}).Error("failed to start transactions")
		return err
	}
	res, err := tx.Exec("UPDATE predictions SET amount=amount+? WHERE playerID=? AND matchID=? AND team=?", amount, discordID, matchID, team)
	if err != nil {
		log.WithFields(log.Fields{
			string(static.UUIDKey):  ctx.Value(static.UUIDKey),
			string(static.ErrorKey): err.Error(),
		}).Error("failed to update predictions")
		tx.Rollback()
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.WithFields(log.Fields{
			string(static.UUIDKey):  ctx.Value(static.UUIDKey),
			string(static.ErrorKey): err.Error(),
		}).Error("failed to get affected rows")
		tx.Rollback()
		return err
	}
	if rowsAffected == 0 {
		_, err := tx.Exec("INSERT INTO predictions (playerID,matchID,team,amount) VALUES (?,?,?,?)", discordID, matchID, team, amount)
		if err != nil {
			log.WithFields(log.Fields{
				string(static.UUIDKey):  ctx.Value(static.UUIDKey),
				string(static.ErrorKey): err.Error(),
			}).Error("failed to insert predictions")
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

func GetPredictionsTotalOnMatch(ctx context.Context, matchID string) (int, int, error) {
	var team1 int
	var team2 int
	row := db.QueryRow("select COALESCE(sum(amount),0) from predictions where matchID=? AND team=1", matchID)
	err := row.Scan(&team1)
	if err != nil {
		return 0, 0, static.ErrDB(err)
	}
	row = db.QueryRow("select COALESCE(sum(amount),0) from predictions where matchID=? AND team=2", matchID)
	err = row.Scan(&team2)
	if err != nil {
		return 0, 0, static.ErrDB(err)
	}
	//we had 50 here to make it for low population
	return team1 + 100, team2 + 100, nil
}
