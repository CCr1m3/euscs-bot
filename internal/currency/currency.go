package currency

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

func GetPlayer(playerID string) (*models.Player, error) {
	return getOrCreatePlayer(playerID)
}

func GetPlayerCurrency(playerID string) (int, error) {
	p, err := getOrCreatePlayer(playerID)
	if err != nil {
		return -1, err
	}
	return p.Currency, nil
}

func AddPrediction(playerID string, matchID string, team int) error {
	_, err := getOrCreatePlayer(playerID)
	if err != nil {
		return err
	}
	err = db.CreatePrediction(playerID, matchID, team)
	if err != nil {
		return err
	}
	return nil
}
