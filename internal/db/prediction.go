package db

import (
	"context"

	"github.com/haashi/omega-strikers-bot/internal/models"
)

func GetPlayersPredictionOnMatch(ctx context.Context, match *models.Match) ([]*models.Prediction, error) {
	predictions := []*models.Prediction{}
	err := db.Select(&predictions, "SELECT elo,discordID,osuser,lastrankupdate,team,credits FROM players JOIN predictions ON predictions.playerID == players.discordID WHERE matchID=? ", match.ID)
	if err != nil {
		return nil, &models.DBError{Err: err}
	}
	return predictions, nil
}

func CreatePrediction(ctx context.Context, discordID string, matchID string, team int) error {
	_, err := db.Exec("INSERT INTO predictions (playerID,matchID,team) VALUES (?,?,?)", discordID, matchID, team)
	if err != nil {
		return &models.DBError{Err: err}
	}
	return nil
}
