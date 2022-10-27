package db

import (
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

func init() {
	logrus.SetLevel(log.DebugLevel)
}

func clearDB() {
	if db != nil {
		err := db.Close()
		if err != nil {
			log.Error("failed to close db: " + err.Error())
		}
	}
	db = nil
	err := os.Remove("omega-strikers-bot.db")
	if err != nil {
		log.Error("error removing file: " + err.Error())
	}
}

func Test_db_getInstance(t *testing.T) {
	clearDB()
	if got := getInstance(); got == nil {
		t.Errorf("getInstance() returned nil")
	}
}

func Test_db_Init(t *testing.T) {
	clearDB()
	Init()
}
