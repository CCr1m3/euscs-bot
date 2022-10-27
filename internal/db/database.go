package db

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
)

var db *sqlx.DB

func getInstance() *sqlx.DB {
	if db == nil {
		var err error
		db, err = sqlx.Open("sqlite3", "./omega-strikers-bot.db")
		if err != nil {
			log.Fatal(err)
		}
		return db
	} else {
		return db
	}
}

func Init() {
	log.Info("starting db service")
	getInstance()
	err := migrate()
	if err != nil {
		log.Fatal(err)
	}
}
