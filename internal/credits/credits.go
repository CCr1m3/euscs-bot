package credits

import (
	"context"

	"github.com/euscs/euscs-bot/internal/db"
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
