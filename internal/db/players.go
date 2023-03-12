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
	if err != nil && !errors.Is(err, static.ErrNotFound) {
		return err
	} else if errors.Is(err, static.ErrNotFound) {
		_, err := db.NamedExec("INSERT INTO players (discordID,twitchID,elo,osuser) VALUES (:discordID,:twitchID,:elo,:osuser)", p)
		if err != nil {
			return static.ErrDB(err)
		}
	} else {
		_, err := db.NamedExec("UPDATE players SET twitchID=:twitchID,elo=:elo,osuser=:osuser WHERE discordID=:discordID", p)
		if err != nil {
			return static.ErrDB(err)
		}
	}
	return nil
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
		p := &Player{DiscordID: discordID}
		err = p.Save(ctx)
		if err != nil {
			return nil, err
		} else {
			return p, nil
		}
	} else if err != nil {
		return nil, err
	}
	return p, nil
}
