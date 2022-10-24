package db

import (
	"time"

	"github.com/haashi/omega-strikers-bot/internal/models"
)

func CreateMatch(m *models.Match) error {
	m.Timestamp = int(time.Now().Unix())
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	_, err = tx.NamedExec("INSERT INTO matches (matchID,threadID,messageID,timestamp) VALUES (:matchID,:threadID,:messageID,:timestamp)", m)
	if err != nil {
		return err
	}
	//add players in matchesplayers
	for _, player := range m.Team1 {
		_, err = tx.Exec("INSERT INTO matchesplayers (matchID,playerID,team) VALUES (?,?,?)", m.ID, player.DiscordID, 1)
		if err != nil {
			return err
		}
	}
	for _, player := range m.Team2 {
		_, err = tx.Exec("INSERT INTO matchesplayers (matchID,playerID,team) VALUES (?,?,?)", m.ID, player.DiscordID, 2)
		if err != nil {
			return err
		}
	}
	err = tx.Commit()
	return err
}

func UpdateMatch(m *models.Match) error {
	//update players in matchesplayers (probably delete and recreate)
	_, err := db.NamedExec("UPDATE matches SET state=:state, team1score=:team1score, team2score=:team2score WHERE matchID=:matchID", m)
	return err
}

func getTeamsInMatch(match *models.Match) error {
	team1 := []*models.Player{}
	err := db.Select(&team1, "SELECT elo,discordID,osuser,lastRankUpdate FROM players JOIN matchesplayers ON matchesplayers.playerID == players.discordID WHERE matchID=? AND team=1", match.ID)
	if err != nil {
		return err
	}
	team2 := []*models.Player{}
	err = db.Select(&team2, "SELECT elo,discordID,osuser,lastRankUpdate FROM players JOIN matchesplayers ON matchesplayers.playerID == players.discordID WHERE matchID=? AND team=2", match.ID)
	if err != nil {
		return err
	}
	match.Team1 = team1
	match.Team2 = team2
	return nil
}

func GetMatchByThreadID(threadID string) (*models.Match, error) {
	var match models.Match
	err := db.Get(&match, "SELECT * FROM matches WHERE threadID=?", threadID)
	if err != nil {
		return nil, err
	}
	err = getTeamsInMatch(&match)
	if err != nil {
		return nil, err
	}
	return &match, nil
}

func GetMatchByID(matchID string) (*models.Match, error) {
	var match models.Match
	err := db.Get(&match, "SELECT * FROM matches WHERE matchID=?", matchID)
	if err != nil {
		return nil, err
	}
	err = getTeamsInMatch(&match)
	if err != nil {
		return nil, err
	}
	return &match, nil
}

func GetRunningMatchesOrderedByTimestamp() ([]*models.Match, error) {
	matches := []*models.Match{}
	err := db.Select(&matches, "SELECT * FROM matches WHERE state=0 ORDER BY timestamp ASC LIMIT 50")
	for _, match := range matches {
		err = getTeamsInMatch(match)
		if err != nil {
			return nil, err
		}
	}
	return matches, err
}

func IsPlayerInMatch(p *models.Player) (bool, error) {
	var count int
	row := db.QueryRow("SELECT COUNT(*) FROM matches JOIN matchesplayers ON matches.matchID = matchesplayers.matchID WHERE playerID=? and state=0", p.DiscordID)
	err := row.Scan(&count)
	return count > 0, err
}
