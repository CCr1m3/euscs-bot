package credits

import (
	"context"

	"github.com/haashi/omega-strikers-bot/internal/db"
)

func GetPlayerCredits(ctx context.Context, playerID string) (int, error) {
	p, err := db.GetOrCreatePlayerById(ctx, playerID)
	if err != nil {
		return -1, err
	}
	return p.Credits, nil
}

func AddPrediction(ctx context.Context, playerID string, matchID string, team int) error {
	_, err := db.GetOrCreatePlayerById(ctx, playerID)
	if err != nil {
		return err
	}
	err = db.CreatePrediction(ctx, playerID, matchID, team)
	if err != nil {
		return err
	}
	return nil
}
