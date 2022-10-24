package rank

import (
	"github.com/haashi/omega-strikers-bot/internal/db"
	"github.com/haashi/omega-strikers-bot/internal/models"
	log "github.com/sirupsen/logrus"
	"os"
	"time"
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
	if err != nil {
		return err
	}
	err = UpdateRankIfNeeded(playerID)
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

func UpdateRankIfNeeded(playerID string) error {
	player, err := db.GetPlayerById(playerID)
	if err != nil {
		return err
	}
	updateDelay := time.Hour * 24
	if os.Getenv("mode") == "dev" {
		updateDelay = time.Second * 5
	}
	if time.Since(time.Unix(int64(player.LastRankUpdate), 0)) > updateDelay {
		log.Infof("updating player elo %s", player.DiscordID)
		rank, err := GetRankFromUsername(player.OSUser)
		if err != nil {
			log.Errorf("Failed to retrieve rank of player %s:"+err.Error(), player.DiscordID)
			return err
		}
		player.Elo = rank
		player.LastRankUpdate = int(time.Now().Unix())
		err = db.UpdatePlayer(player)
		if err != nil {
			log.Errorf("Failed to update player %s:"+err.Error(), player.DiscordID)
		}
		return err
	}
	return nil
}
