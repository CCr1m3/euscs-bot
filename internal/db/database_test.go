package db

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/euscs/euscs-bot/internal/env"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
)

func init() {
	err := godotenv.Load("../../.env")
	env.Init()
	if err != nil {
		log.Warning("error loading .env file: " + err.Error())
	}
	log.SetOutput(ioutil.Discard)
}

func clearDB() {
	getInstance()
	if db != nil {
		if env.DB.Type == env.SQLITE {
			err := db.Close()
			if err != nil {
				log.Error("failed to close db: " + err.Error())
			}
			err = os.Remove(env.DB.Path)
			if err != nil {
				log.Error("error removing file: " + err.Error())
			}
		}
		if env.DB.Type == env.MYSQL {
			tx, err := db.Beginx()
			if err != nil {
				log.Error("error starting transaction:" + err.Error())
			}
			rows := []struct {
				Tables_in_euos string `db:"Tables_in_euos"`
			}{}
			tx.Exec("SET foreign_key_checks = 0")
			err = tx.Select(&rows, "SHOW TABLES in euos")
			if err != nil {
				log.Error("error getting database table:" + err.Error())
			}
			for _, row := range rows {
				_, err := tx.Exec("DROP TABLE " + row.Tables_in_euos)
				if err != nil {
					log.Error("error dropping table:" + err.Error())
				}
			}
			tx.Exec("SET foreign_key_checks = 1")
			tx.Commit()
			err = db.Close()
			if err != nil {
				log.Error("failed to close db: " + err.Error())
			}
		}
	}
	db = nil
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
