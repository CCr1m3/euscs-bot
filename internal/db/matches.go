package db

import "time"

type Match struct {
	Team1      []*Player
	Team2      []*Player
	ThreadID   string `db:"threadID"`
	MessageID  string `db:"messageID"`
	ID         string `db:"matchID"`
	Running    int    `db:"running"`
	Team1Score int    `db:"team1score"`
	Team2Score int    `db:"team2score"`
	Timestamp  int    `db:"timestamp"`
}

func (m *Match) Save() error {
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

func (m *Match) SetScore(team1Score int, team2Score int) error {
	m.Team1Score = team1Score
	m.Team2Score = team2Score
	m.Running = 0
	_, err := db.NamedExec("UPDATE matches SET running=:running, team1score=:team1score, team2score=:team2score WHERE matchID=:matchID", m)
	return err
}

func GetMatchByThreadID(threadID string) (*Match, error) {
	var match Match
	err := db.Get(&match, "SELECT * FROM matches WHERE threadID=?", threadID)
	if err != nil {
		return nil, err
	}
	return &match, nil
}

func GetMatchByID(matchID string) (*Match, error) {
	var match Match
	err := db.Get(&match, "SELECT * FROM matches WHERE matchID=?", matchID)
	if err != nil {
		return nil, err
	}
	return &match, nil
}

func GetMatchesByTimestamp() ([]*Match, error) {
	matches := []*Match{}
	err := db.Select(&matches, "SELECT * FROM matches WHERE running=1 ORDER BY timestamp ASC LIMIT 50")
	return matches, err
}
