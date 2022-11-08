package db

import (
	"context"

	"github.com/haashi/omega-strikers-bot/internal/models"
)

func AddPlayerToQueue(ctx context.Context, p *models.Player, role models.Role, entryTime int) error {
	_, err := db.Exec("INSERT INTO queue (playerID,role,entryTime) VALUES (?,?,?)", p.DiscordID, role, entryTime)
	if err != nil {
		return &models.DBError{Err: err}
	}
	return nil
}

func RemovePlayerFromQueue(ctx context.Context, p *models.Player) error {
	_, err := db.NamedExec("DELETE FROM queue WHERE playerID=:discordID", p)
	if err != nil {
		return &models.DBError{Err: err}
	}
	return nil
}

func GetPlayersInQueue(ctx context.Context) ([]*models.QueuedPlayer, error) {
	players := []*models.QueuedPlayer{}
	err := db.Select(&players, "SELECT discordID,osuser,elo,role,lastrankupdate,credits,entrytime FROM queue JOIN players ON queue.playerID = players.discordID")
	if err != nil {
		return nil, &models.DBError{Err: err}
	}
	return players, nil
}

func IsPlayerInQueue(ctx context.Context, p *models.Player) (bool, error) {
	var count int
	row := db.QueryRow("SELECT COUNT(*) FROM queue WHERE playerID=?", p.DiscordID)
	err := row.Scan(&count)
	if err != nil {
		return false, &models.DBError{Err: err}
	}
	return count > 0, nil
}

func GetGoaliesCountInQueue(ctx context.Context) (int, error) {
	var count int
	row := db.QueryRow("SELECT COUNT(*) FROM queue WHERE role='goalie' OR role='flex'")
	err := row.Scan(&count)
	if err != nil {
		return 0, &models.DBError{Err: err}
	}
	return count, nil
}

func GetForwardsCountInQueue(ctx context.Context) (int, error) {
	var count int
	row := db.QueryRow("SELECT COUNT(*) FROM queue WHERE role='forward' OR role='flex'")
	err := row.Scan(&count)
	if err != nil {
		return 0, &models.DBError{Err: err}
	}
	return count, nil
}
