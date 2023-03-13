package db

import (
	"context"
	"database/sql"
	"errors"

	"github.com/euscs/euscs-bot/internal/static"
)

type Player struct {
	DiscordID string `db:"discordID"`
	TwitchID  string `db:"twitchID"`
	Elo       int    `db:"elo"`
	OSUser    string `db:"osuser"`
}
type Players []*Player

func (p *Player) Save(ctx context.Context) error {
	if p.DiscordID == "" {
		return static.ErrDiscordIDRequired
	}
	_, err := GetPlayerByID(ctx, p.DiscordID)
	if err != nil {
		return err
	}
	_, err = db.NamedExec("UPDATE players SET twitchID=:twitchID,elo=:elo,osuser=:osuser WHERE discordID=:discordID", p)
	if err != nil {
		return static.ErrDB(err)
	}
	return nil
}

func CreatePlayerWithID(ctx context.Context, discordID string) (*Player, error) {
	_, err := GetPlayerByID(ctx, discordID)
	if err != nil && !errors.Is(err, static.ErrNotFound) {
		return nil, err
	} else if err == nil {
		return nil, static.ErrAlreadyExists
	}
	_, err = db.Exec("INSERT INTO players (discordID) VALUES (?)", discordID)
	if err != nil {
		return nil, static.ErrDB(err)
	}
	return &Player{DiscordID: discordID, Elo: 1500}, nil
}

func GetPlayerByID(ctx context.Context, discordID string) (*Player, error) {
	var player Player
	err := db.Get(&player, "SELECT discordID,twitchID,elo,osuser FROM players WHERE discordID=?", discordID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, static.ErrNotFound
		}
		return nil, static.ErrDB(err)
	}
	return &player, nil
}

func GetPlayerByUsername(ctx context.Context, osuser string) (*Player, error) {
	var player Player
	err := db.Get(&player, "SELECT discordID,twitchID,elo,osuser FROM players WHERE osuser=?", osuser)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, static.ErrNotFound
		}
		return nil, static.ErrDB(err)
	}
	return &player, nil
}

func GetOrCreatePlayerByID(ctx context.Context, discordID string) (*Player, error) {
	p, err := GetPlayerByID(ctx, discordID)
	if err != nil && errors.Is(err, static.ErrNotFound) {
		return CreatePlayerWithID(ctx, discordID)
	} else if err != nil {
		return nil, err
	}
	return p, nil
}

func (p *Player) isDummy() bool {
	return len(p.DiscordID) < 10
}
