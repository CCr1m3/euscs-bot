package db

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/euscs/euscs-bot/internal/static"
)

type Player struct {
	DiscordID      string `db:"discordID"`
	Elo            int    `db:"elo"`
	OSUser         string `db:"osuser"`
	LastRankUpdate int    `db:"lastrankupdate"`
	Credits        int    `db:"credits"`
}
type Players []*Player

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
	_, err = db.Exec("UPDATE players SET lastrankupdate=? WHERE discordID=?", time.Now(), p.DiscordID)
	if err != nil {
		return static.ErrDB(err)
	}
	p.Elo = elo
	return nil
}

func (p *Player) SetCredits(ctx context.Context, credits int) error {
	_, err := GetPlayerByID(ctx, p.DiscordID)
	if err != nil {
		return err
	}
	_, err = db.Exec("UPDATE players SET credits=? WHERE discordID=?", credits, p.DiscordID)
	if err != nil {
		return static.ErrDB(err)
	}
	p.Credits = credits
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
	err := db.Get(&player, "SELECT discordID,elo,osuser,lastrankupdate,credits FROM players WHERE discordID=?", discordID)
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
	err := db.Get(&player, "SELECT discordID,elo,osuser,lastrankupdate,credits FROM players WHERE osuser=?", osuser)
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
