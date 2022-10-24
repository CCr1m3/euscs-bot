package db

import "github.com/haashi/omega-strikers-bot/internal/models"

func AddPlayerToQueue(p *models.Player, role models.Role) error {
	_, err := db.Exec("INSERT INTO queue (playerID,role) VALUES (?,?)", p.DiscordID, role)
	return err
}

func RemovePlayerFromQueue(p *models.Player) error {
	_, err := db.NamedExec("DELETE FROM queue WHERE playerID=:discordID", p)
	return err
}

func GetPlayersInQueue() ([]*models.QueuedPlayer, error) {
	players := []*models.QueuedPlayer{}
	err := db.Select(&players, "SELECT discordID,osuser,elo,role FROM queue JOIN players ON queue.playerID = players.discordID")
	return players, err
}

func IsPlayerInQueue(p *models.Player) (bool, error) {
	var count int
	row := db.QueryRow("SELECT COUNT(*) FROM queue WHERE playerID=?", p.DiscordID)
	err := row.Scan(&count)
	return count > 0, err
}

func GetGoaliesCountInQueue() (int, error) {
	var count int
	row := db.QueryRow("SELECT COUNT(*) FROM queue WHERE role='goalie' OR role='flex'")
	err := row.Scan(&count)
	return count, err
}

func GetForwardsCountInQueue() (int, error) {
	var count int
	row := db.QueryRow("SELECT COUNT(*) FROM queue WHERE (role='forward' OR role='flex')")
	err := row.Scan(&count)
	return count, err
}
