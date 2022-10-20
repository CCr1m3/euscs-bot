package db

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
)

var instance *sqlx.DB
var db *sqlx.DB

func getInstance() *sqlx.DB {
	if instance == nil {
		db, err := sqlx.Open("sqlite3", "./omega-strikers-bot.db")
		if err != nil {
			log.Fatal(err)
		}
		instance = db
		return instance
	} else {
		return instance
	}
}

func Init() {
	log.Info("starting db service")
	db = getInstance()
	err := migrate()
	if err != nil {
		log.Fatal(err)
	}
}
