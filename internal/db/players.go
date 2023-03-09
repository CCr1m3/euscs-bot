package db

import (
	"context"

	"github.com/euscs/euscs-bot/internal/models"
)

func CreatePlayer(ctx context.Context, discordID string) error {
	_, err := db.Exec("INSERT INTO players (discordID) VALUES (?)", discordID)
	if err != nil {
		return &models.DBError{Err: err}
	}
	return nil
}

func GetPlayerById(ctx context.Context, discordID string) (*models.Player, error) {
	var player models.Player
	err := db.Get(&player, "SELECT * FROM players WHERE discordID=?", discordID)
	if err != nil {
		return nil, &models.DBError{Err: err}
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
		return nil, &models.DBError{Err: err}
	}
	return &player, nil
}

func UpdatePlayer(ctx context.Context, p *models.Player) error {
	_, err := db.NamedExec("UPDATE players SET twitchID=:twitchID,elo=:elo,osuser=:osuser,lastrankupdate=:lastrankupdate,credits=:credits WHERE discordID=:discordID", p)
	if err != nil {
		return &models.DBError{Err: err}
	}
	return nil
}

func GetPlayersOrderedByCredits(ctx context.Context) ([]*models.Player, error) {
	predictions := []*models.Player{}
	err := db.Select(&predictions, "SELECT * FROM players ORDER BY credits DESC")
	if err != nil {
		return nil, &models.DBError{Err: err}
	}
	return predictions, nil
}
