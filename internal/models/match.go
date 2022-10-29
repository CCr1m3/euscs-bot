package models

type Match struct {
	Team1         []*Player
	Team2         []*Player
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
	MatchStateInProgress     MatchState = 0
	MatchStateTeam1Won       MatchState = 1
	MatchStateTeam2Won       MatchState = 2
	MatchStateVoteInProgress MatchState = -1
	MatchStateCanceled       MatchState = 3
)

type Prediction struct {
	Player
	Team int `db:"team"`
}
