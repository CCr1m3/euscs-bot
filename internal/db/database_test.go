package db

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Warning("error loading .env file: " + err.Error())
	}
	logLevel := os.Getenv("loglevel")
	if logLevel == "debug" {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}
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
