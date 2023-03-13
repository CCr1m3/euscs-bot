package db

import (
	"io/ioutil"
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

func Test_db_getInstance(t *testing.T) {
	Clear()
	if got := GetInstance(); got == nil {
		t.Errorf("getInstance() returned nil")
	}
}

func Test_db_Init(t *testing.T) {
	Clear()
	Init()
}
