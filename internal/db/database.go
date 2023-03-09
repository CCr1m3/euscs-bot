package db

import (
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/haashi/omega-strikers-bot/internal/env"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
)

var db *sqlx.DB

func getInstance() *sqlx.DB {
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
	getInstance()
	err := migrate()
	if err != nil {
		log.Fatal(err)
	}
}
