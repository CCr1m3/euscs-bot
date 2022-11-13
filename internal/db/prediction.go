package db

import (
	"context"

	"github.com/haashi/omega-strikers-bot/internal/models"
)

func GetPlayersPredictionOnMatch(ctx context.Context, match *models.Match) ([]*models.Prediction, error) {
	predictions := []*models.Prediction{}
	err := db.Select(&predictions, "SELECT elo,discordID,osuser,lastrankupdate,team,credits,amount FROM players JOIN predictions ON predictions.playerID == players.discordID WHERE matchID=? ", match.ID)
	if err != nil {
		return nil, &models.DBError{Err: err}
	}
	return predictions, nil
}

func CreatePrediction(ctx context.Context, discordID string, matchID string, team int, amount int) error {
	_, err := db.Exec("INSERT INTO predictions (playerID,matchID,team,amount) VALUES (?,?,?,?)", discordID, matchID, team, amount)
	if err != nil {
		return &models.DBError{Err: err}
	}
	return nil
}
