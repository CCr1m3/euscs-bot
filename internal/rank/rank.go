package rank

import (
	"github.com/haashi/omega-strikers-bot/internal/db"
	"github.com/haashi/omega-strikers-bot/internal/models"
	log "github.com/sirupsen/logrus"
)

func Init() {
	log.Info("starting rank service")
}

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

func LinkPlayerToUsername(playerID string, username string) error {
	player, err := getOrCreatePlayer(playerID)
	if err != nil {
		return err
	}
	player.OSUser = username
	err = db.UpdatePlayer(player)
	return err
}

func GetLinkedUsername(playerID string) (string, error) {
	player, err := getOrCreatePlayer(playerID)
	if err != nil {
		return "", err
	}
	return player.OSUser, nil
}

func GetLinkedUser(username string) (string, error) {
	player, err := db.GetPlayerByUsername(username)
	if err != nil {
		return "", err
	}
	return player.DiscordID, nil
}
