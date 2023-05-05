package db

import (
	"context"

	"github.com/euscs/euscs-bot/internal/static"
)

type Role string

const (
	RoleFlex    Role = "flex"
	RoleForward Role = "forward"
	RoleGoalie  Role = "goalie"
)

type QueuedPlayer struct {
	Player
	Role      Role `db:"role"`
	EntryTime int  `db:"entrytime"`
}

func AddPlayerToQueue(ctx context.Context, p *Player, role Role, entryTime int) error {
	_, err := db.Exec("INSERT INTO queue (playerID,role,entryTime) VALUES (?,?,?)", p.DiscordID, role, entryTime)
	if err != nil {
		return static.ErrDB(err)
	}
	return nil
}

func RemovePlayerFromQueue(ctx context.Context, p *Player) error {
	_, err := db.NamedExec("DELETE FROM queue WHERE playerID=:discordID", p)
	if err != nil {
		return static.ErrDB(err)
	}
	return nil
}

func GetPlayersInQueue(ctx context.Context) ([]*QueuedPlayer, error) {
	players := []*QueuedPlayer{}
	err := db.Select(&players, "SELECT discordID,osuser,elo,role,lastrankupdate,credits,entrytime FROM queue JOIN players ON queue.playerID = players.discordID")
	if err != nil {
		return nil, static.ErrDB(err)
	}
	return players, nil
}

func IsPlayerInQueue(ctx context.Context, p *Player) (bool, error) {
	var count int
	row := db.QueryRow("SELECT COUNT(*) FROM queue WHERE playerID=?", p.DiscordID)
	err := row.Scan(&count)
	if err != nil {
		return false, static.ErrDB(err)
	}
	return count > 0, nil
}

func GetGoaliesCountInQueue(ctx context.Context) (int, error) {
	var count int
	row := db.QueryRow("SELECT COUNT(*) FROM queue WHERE role='goalie' OR role='flex'")
	err := row.Scan(&count)
	if err != nil {
		return 0, static.ErrDB(err)
	}
	return count, nil
}

func GetForwardsCountInQueue(ctx context.Context) (int, error) {
	var count int
	row := db.QueryRow("SELECT COUNT(*) FROM queue WHERE role='forward' OR role='flex'")
	err := row.Scan(&count)
	if err != nil {
		return 0, static.ErrDB(err)
	}
	return count, nil
}
