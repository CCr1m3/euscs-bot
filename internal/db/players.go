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

func (p *Player) SetTwitchID(ctx context.Context, twitchID string) error {
	_, err := GetPlayerByID(ctx, p.DiscordID)
	if err != nil {
		return err
	}
	_, err = db.Exec("UPDATE players SET twitchID=? WHERE discordID=?", twitchID, p.DiscordID)
	if err != nil {
		return static.ErrDB(err)
	}
	p.TwitchID = twitchID
	return nil
}

func (p *Player) SetOSUser(ctx context.Context, OSUser string) error {
	_, err := GetPlayerByID(ctx, p.DiscordID)
	if err != nil {
		return err
	}
	_, err = db.Exec("UPDATE players SET osuser=? WHERE discordID=?", OSUser, p.DiscordID)
	if err != nil {
		return static.ErrDB(err)
	}
	p.OSUser = OSUser
	return nil
}

func (p *Player) SetElo(ctx context.Context, elo int) error {
	_, err := GetPlayerByID(ctx, p.DiscordID)
	if err != nil {
		return err
	}
	_, err = db.Exec("UPDATE players SET elo=? WHERE discordID=?", elo, p.DiscordID)
	if err != nil {
		return static.ErrDB(err)
	}
	p.Elo = elo
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
