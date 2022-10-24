package matchmaking

import (
	"github.com/haashi/omega-strikers-bot/internal/db"
	"github.com/haashi/omega-strikers-bot/internal/models"
)

func getOrCreatePlayer(playerID string) (*models.Player, error) {
	p, err := db.GetPlayerById(playerID)
	if err != nil {
		err = db.CreatePlayer(playerID)
		if err != nil {
			return nil, err
		}
		p, err = db.GetPlayerById(playerID)
		if err != nil {
			return nil, err
		}
	}
	return p, nil
}
