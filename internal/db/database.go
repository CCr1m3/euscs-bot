package db

import (
	"os"
	"time"

	"github.com/euscs/euscs-bot/internal/env"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
)

var db *sqlx.DB

func GetInstance() *sqlx.DB {
	if db == nil {
		var err error
		if env.DB.Type == env.SQLITE {
			db, err = sqlx.Open("sqlite3", env.DB.Path)
			if err != nil {
				log.Fatal(err)
			}
		}
		if env.DB.Type == env.MYSQL {
			db, err = sqlx.Open("mysql", env.DB.Path)
			if err != nil {
				log.Fatal(err)
			}
			db.SetConnMaxLifetime(time.Minute * 5)
			db.SetMaxOpenConns(1)
			db.SetMaxIdleConns(1)
		}
		return db
	} else {
		return db
	}
}

func Init() {
	log.Info("starting db service")
	GetInstance()
	err := migrate()
	if err != nil {
		log.Fatal(err)
	}
}

func Clear() {
	GetInstance()
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
