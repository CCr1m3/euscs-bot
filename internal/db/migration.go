package db

import (
	"fmt"
	"strings"

	"github.com/euscs/euscs-bot/internal/models"
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
	if err != nil && !strings.Contains(err.Error(), "UNIQUE") && !strings.Contains(err.Error(), "1062") {
		return models.ErrDB(err)
	}
	start, err = getLatestMigration()
	if err != nil {
		return models.ErrDB(err)
	}
	for i := start + 1; i < len(migrations); i++ {
		log.Info(fmt.Sprintf("applying migration %d", i))
		_, err = db.Exec(migrations[i])
		if err != nil {
			return models.ErrDB(err)
		}
		_, err = db.Exec(`INSERT INTO migrations (version) VALUES (?)`, i)
		if err != nil {
			return models.ErrDB(err)
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
		return 0, models.ErrDB(err)
	}
	return ver, err
}

var migration0 = `CREATE TABLE IF NOT EXISTS migrations (version INTEGER,PRIMARY KEY (version));
INSERT INTO migrations (version) VALUES (0);`

var migration1 = `CREATE TABLE players (
    discordID VARCHAR(100) UNIQUE NOT NULL,
		elo INTEGER DEFAULT 1500 NOT NULL,
		osuser VARCHAR(100) DEFAULT "" NOT NULL,
		twitchID VARCHAR(100) default "",
		PRIMARY KEY (discordID)
);`

var migration2 = `CREATE TABLE teams (
	name VARCHAR(100) UNIQUE NOT NULL,
	ownerplayerID VARCHAR(100) UNIQUE NOT NULL,
	FOREIGN KEY (ownerplayerID) REFERENCES players(discordID),
	PRIMARY KEY(name)
);
CREATE TABLE teamsplayers (
	playerID VARCHAR(100) UNIQUE NOT NULL,
	team VARCHAR(100) NOT NULL,
	FOREIGN KEY (playerID) REFERENCES players(discordID),
	FOREIGN KEY (team) REFERENCES teams(name),
	PRIMARY KEY (playerID,team)
);`

var migration3 = `CREATE TABLE markov (
	word1 VARCHAR(100) NOT NULL,
	word2	VARCHAR(100) NOT NULL,
	word3	VARCHAR(100) NOT NULL,
	count INTEGER NOT NULL,
	PRIMARY KEY (word1,word2,word3)
);`

/*var migration4 = `CREATE TABLE matches (
	matchID VARCHAR(100) UNIQUE NOT NULL,
	messageID VARCHAR(100) UNIQUE,
	votemessageID VARCHAR(100),
	threadID VARCHAR(100) UNIQUE,
	timestamp INTEGER NOT NULL,
	state INTEGER DEFAULT 0 NOT NULL,
	team1score INTEGER DEFAULT 0 NOT NULL,
	team2score INTEGER DEFAULT 0 NOT NULL,
	PRIMARY KEY(matchID)
);
CREATE TABLE matchesplayers (
	matchID VARCHAR(100) NOT NULL,
	team INTEGER NOT NULL,
	playerID VARCHAR(100) NOT NULL,
	FOREIGN KEY (playerID) REFERENCES players(discordID),
	FOREIGN KEY (matchID) REFERENCES matches(matchID),
	PRIMARY KEY (playerID,matchID)
);`*/
