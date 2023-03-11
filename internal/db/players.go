package db

import (
	"context"
	"database/sql"

	"github.com/euscs/euscs-bot/internal/models"
)

func CreatePlayer(ctx context.Context, discordID string) error {
	_, err := db.Exec("INSERT INTO players (discordID) VALUES (?)", discordID)
	if err != nil {
		return models.ErrDB(err)
	}
	return nil
}

func GetPlayerById(ctx context.Context, discordID string) (*models.Player, error) {
	var player models.Player
	err := db.Get(&player, "SELECT discordID,twitchID,elo,osuser FROM players WHERE discordID=?", discordID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, models.ErrNotFound
		}
		return nil, models.ErrDB(err)
	}
	return &player, nil
}

func GetOrCreatePlayerById(ctx context.Context, discordID string) (*models.Player, error) {
	p, err := GetPlayerById(ctx, discordID)
	if err != nil {
		err = CreatePlayer(ctx, discordID)
		if err != nil {
			return nil, err
		}
		p, err = GetPlayerById(ctx, discordID)
		if err != nil {
			return nil, err
		}
	}
	return p, nil
}

func GetPlayerByUsername(ctx context.Context, username string) (*models.Player, error) {
	var player models.Player
	err := db.Get(&player, "SELECT * FROM players WHERE osuser=?", username)
	if err != nil {
		return nil, models.ErrDB(err)
	}
	return &player, nil
}

func UpdatePlayer(ctx context.Context, p *models.Player) error {
	_, err := db.NamedExec("UPDATE players SET twitchID=:twitchID,elo=:elo,osuser=:osuser WHERE discordID=:discordID", p)
	if err != nil {
		return models.ErrDB(err)
	}
	return nil
}
