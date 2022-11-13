package credits

import (
	"context"

	"github.com/haashi/omega-strikers-bot/internal/db"
	"github.com/haashi/omega-strikers-bot/internal/models"
	log "github.com/sirupsen/logrus"
)

func GetPlayerCredits(ctx context.Context, playerID string) (int, error) {
	p, err := db.GetOrCreatePlayerById(ctx, playerID)
	if err != nil {
		return -1, err
	}
	return p.Credits, nil
}

func AddPrediction(ctx context.Context, playerID string, matchID string, team int, amount int) error {
	player, err := db.GetOrCreatePlayerById(ctx, playerID)
	if err != nil {
		return err
	}
	err = db.CreatePrediction(ctx, playerID, matchID, team, amount)
	if err != nil {
		return err
	}
	player.Credits -= amount
	err = db.UpdatePlayer(ctx, player)
	if err != nil {
		return err
	}
	return nil
}

func GetReturnRatiosForMatch(ctx context.Context, matchID string) (float64, float64, error) {
	match, err := db.GetMatchByID(ctx, matchID)
	if err != nil {
		log.WithFields(log.Fields{
			string(models.UUIDKey):    ctx.Value(models.UUIDKey),
			string(models.ErrorKey):   err.Error(),
			string(models.MatchIDKey): matchID,
		}).Error("failed to get match")
		return 0, 0, err
	}
	predictions, err := db.GetPlayersPredictionOnMatch(ctx, match)
	if err != nil {
		log.WithFields(log.Fields{
			string(models.UUIDKey):    ctx.Value(models.UUIDKey),
			string(models.ErrorKey):   err.Error(),
			string(models.MatchIDKey): matchID,
		}).Error("failed to get match predictions")
		return 0, 0, err
	}
	team1sum := 0
	team2sum := 0
	for _, prediction := range predictions {
		if prediction.Team == 1 {
			team1sum += prediction.Amount
		} else if prediction.Team == 2 {
			team2sum += prediction.Amount
		}
	}
	if team1sum == 0 {
		team1sum++
	}
	if team2sum == 0 {
		team2sum++
	}
	ratioTeam1 := float64(team2sum) / float64(team1sum)
	ratioTeam2 := 1 / ratioTeam1
	return ratioTeam1 + 1, ratioTeam2 + 1, nil
}
