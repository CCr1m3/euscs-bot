package db

import (
	"context"

	"github.com/euscs/euscs-bot/internal/static"
)

type Match struct {
	Team1         []*Player  `db:"team1"`
	Team2         []*Player  `db:"team2"`
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
	MatchStateVoteInProgress MatchState = -1
	MatchStateInProgress     MatchState = 0
	MatchStateTeam1Won       MatchState = 1
	MatchStateTeam2Won       MatchState = 2
	MatchStateCanceled       MatchState = 3
)

func CreateMatch(ctx context.Context, m *Match) error {
	tx, err := db.Beginx()
	if err != nil {
		return static.ErrDB(err)
	}
	_, err = tx.NamedExec("INSERT INTO matches (matchID,threadID,messageID,votemessageID,timestamp) VALUES (:matchID,:threadID,:messageID,:votemessageID,:timestamp)", m)
	if err != nil {
		err2 := tx.Rollback()
		if err2 != nil {
			return static.ErrDB(err2)
		}
		return static.ErrDB(err)
	}
	for _, player := range m.Team1 {
		_, err = tx.Exec("INSERT INTO matchesplayers (matchID,playerID,team) VALUES (?,?,?)", m.ID, player.DiscordID, 1)
		if err != nil {
			err2 := tx.Rollback()
			if err2 != nil {
				return static.ErrDB(err2)
			}
			return static.ErrDB(err)
		}
	}
	for _, player := range m.Team2 {
		_, err = tx.Exec("INSERT INTO matchesplayers (matchID,playerID,team) VALUES (?,?,?)", m.ID, player.DiscordID, 2)
		if err != nil {
			err2 := tx.Rollback()
			if err2 != nil {
				return static.ErrDB(err2)
			}
			return static.ErrDB(err)
		}
	}
	err = tx.Commit()
	if err != nil {
		err2 := tx.Rollback()
		if err2 != nil {
			return static.ErrDB(err2)
		}
		return static.ErrDB(err)
	}
	return nil
}

func UpdateMatch(ctx context.Context, m *Match) error {
	//update players in matchesplayers (probably delete and recreate)
	_, err := db.NamedExec("UPDATE matches SET state=:state, team1score=:team1score, team2score=:team2score,votemessageID=:votemessageID WHERE matchID=:matchID", m)
	if err != nil {
		return static.ErrDB(err)
	}
	return nil
}

func getTeamsInMatch(ctx context.Context, match *Match) error {
	team1 := []*Player{}
	err := db.Select(&team1, "SELECT elo,discordID,osuser,lastrankupdate,credits FROM players JOIN matchesplayers ON matchesplayers.playerID = players.discordID WHERE matchID=? AND team=1", match.ID)
	if err != nil {
		return static.ErrDB(err)
	}
	team2 := []*Player{}
	err = db.Select(&team2, "SELECT elo,discordID,osuser,lastrankupdate,credits FROM players JOIN matchesplayers ON matchesplayers.playerID = players.discordID WHERE matchID=? AND team=2", match.ID)
	if err != nil {
		return static.ErrDB(err)
	}
	match.Team1 = team1
	match.Team2 = team2
	return nil
}

func GetMatchByThreadID(ctx context.Context, threadID string) (*Match, error) {
	var match Match
	err := db.Get(&match, "SELECT * FROM matches WHERE threadID=?", threadID)
	if err != nil {
		return nil, static.ErrDB(err)
	}
	err = getTeamsInMatch(ctx, &match)
	if err != nil {
		return nil, static.ErrDB(err)
	}
	return &match, nil
}

func GetMatchByID(ctx context.Context, matchID string) (*Match, error) {
	var match Match
	err := db.Get(&match, "SELECT * FROM matches WHERE matchID=?", matchID)
	if err != nil {
		return nil, static.ErrDB(err)
	}
	err = getTeamsInMatch(ctx, &match)
	if err != nil {
		return nil, static.ErrDB(err)
	}
	return &match, nil
}

func GetRunningMatchesOrderedByTimestamp(ctx context.Context) ([]*Match, error) {
	matches := []*Match{}
	err := db.Select(&matches, "SELECT * FROM matches WHERE state=0 ORDER BY timestamp ASC LIMIT 50")
	if err != nil {
		return nil, static.ErrDB(err)
	}
	for _, match := range matches {
		err = getTeamsInMatch(ctx, match)
		if err != nil {
			return nil, static.ErrDB(err)
		}
	}
	return matches, nil
}

func GetWaitingForVotesMatches(ctx context.Context) ([]*Match, error) {
	matches := []*Match{}
	err := db.Select(&matches, "SELECT * FROM matches WHERE state=-1")
	if err != nil {
		return nil, static.ErrDB(err)
	}
	for _, match := range matches {
		err = getTeamsInMatch(ctx, match)
		if err != nil {
			return nil, static.ErrDB(err)
		}
	}
	return matches, nil
}

func IsPlayerInMatch(ctx context.Context, p *Player) (bool, error) {
	var count int
	row := db.QueryRow("SELECT COUNT(*) FROM matches JOIN matchesplayers ON matches.matchID = matchesplayers.matchID WHERE playerID=? and state<=0", p.DiscordID)
	err := row.Scan(&count)
	if err != nil {
		return false, static.ErrDB(err)
	}
	return count > 0, nil
}
