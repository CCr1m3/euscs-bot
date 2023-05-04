package db

type Match struct {
	Team1         []Player   `db:"team1"`
	Team2         []Player   `db:"team2"`
	ThreadID      string     `db:"threadID"`
	MessageID     string     `db:"messageID"`
	VoteMessageID string     `db:"votemessageID"`
	ID            string     `db:"matchID"`
	State         MatchState `db:"state"`
	Team1Score    int        `db:"team1score"`
	Team2Score    int        `db:"team2score"`
	Timestamp     int        `db:"timestamp"`
}

type MatchState int

const (
	MatchStateInProgress MatchState = 0
	MatchStateTeam1Won   MatchState = 1
	MatchStateTeam2Won   MatchState = 2
	MatchStateCanceled   MatchState = 3
)
