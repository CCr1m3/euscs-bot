package db

import (
	"github.com/haashi/omega-strikers-bot/internal/models"
)

func CreatePlayer(discordID string) error {
	_, err := db.Exec("INSERT INTO players (discordID) VALUES (?)", discordID)
	return err
}

func GetPlayerById(discordID string) (*models.Player, error) {
	var player models.Player
	err := db.Get(&player, "SELECT * FROM players WHERE discordID=?", discordID)
	if err != nil {
		return nil, err
	}
	return &player, nil
}

func GetPlayerByUsername(username string) (*models.Player, error) {
	var player models.Player
	err := db.Get(&player, "SELECT * FROM players WHERE osuser=?", username)
	if err != nil {
		return nil, err
	}
	return &player, nil
}

func UpdatePlayer(p *models.Player) error {
	_, err := db.NamedExec("UPDATE players SET elo=:elo,osuser=:osuser, lastRankUpdate=:lastRankUpdate WHERE discordID=:discordID", p)
	return err
}
