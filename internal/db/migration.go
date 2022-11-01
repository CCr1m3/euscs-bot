package db

import (
	"fmt"
	"strings"

	"github.com/haashi/omega-strikers-bot/internal/models"
	log "github.com/sirupsen/logrus"
)

var migrations = []string{
	migration0,
	migration1,
	migration2,
	migration3,
}

func migrate() error {
	var start int
	_, err := db.Exec(migrations[0])
	if err != nil && !strings.Contains(err.Error(), "UNIQUE") {
		return &models.DBError{Err: err}
	}
	start, err = getLatestMigration()
	if err != nil {
		return &models.DBError{Err: err}
	}
	for i := start + 1; i < len(migrations); i++ {
		log.Info(fmt.Sprintf("applying migration %d", i))
		_, err = db.Exec(migrations[i])
		if err != nil {
			return &models.DBError{Err: err}
		}
		_, err = db.Exec(`INSERT INTO migrations (version) VALUES (?)`, i)
		if err != nil {
			return &models.DBError{Err: err}
		}
	}
	return nil
}

func getLatestMigration() (int, error) {
	ver := 0
	row := db.QueryRow(`SELECT version
	FROM migrations
	ORDER BY version DESC
	LIMIT 1`)
	err := row.Scan(&ver)
	if err != nil {
		return 0, &models.DBError{Err: err}
	}
	return ver, err
}

var migration0 = `CREATE TABLE IF NOT EXISTS migrations (
	version int,
	PRIMARY KEY (version)
);
INSERT INTO migrations (version) VALUES (0);
`

var migration1 = `CREATE TABLE players (
    discordID text UNIQUE,
		elo int DEFAULT 1500 NOT NULL,
		osuser text DEFAULT "",
		lastrankupdate int DEFAULT 0 NOT NULL,
		credits int DEFAULT 0 NOT NULL,
		PRIMARY KEY (discordID)
);
CREATE TABLE queue (
	playerID text UNIQUE,
	role text DEFAULT "" NOT NULL,
	entrytime int NOT NULL,
	PRIMARY KEY (playerID),
	FOREIGN KEY (playerID) REFERENCES players(discordID)
);`

var migration2 = `CREATE TABLE matches (
	matchID text UNIQUE,
	messageID text UNIQUE,
	votemessageID text DEFAULT "",
	threadID text UNIQUE,
	timestamp int,
	state int DEFAULT 0 NOT NULL,
	team1score int DEFAULT 0 NOT NULL,
	team2score int DEFAULT 0 NOT NULL,
	PRIMARY KEY(matchID)
);
CREATE TABLE matchesplayers (
	matchID text,
	team int,
	playerID text,
	FOREIGN KEY (playerID) REFERENCES players(discordID),
	FOREIGN KEY (matchID) REFERENCES matches(matchID),
	PRIMARY KEY (playerID,matchID)
);
CREATE TABLE predictions (
	matchID text,
	team int,
	playerID text,
	FOREIGN KEY (playerID) REFERENCES players(discordID),
	FOREIGN KEY (matchID) REFERENCES matches(matchID),
	PRIMARY KEY (playerID,matchID)
)
`
var migration3 = `CREATE TABLE markov (
	word1 text,
	word2	text,
	word3	text,
	count int,
	PRIMARY KEY (word1,word2,word3)
);
`
