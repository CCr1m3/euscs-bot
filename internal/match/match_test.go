package match

import (
	"context"
	"io"
	"testing"

	"github.com/euscs/euscs-bot/internal/db"
	"github.com/euscs/euscs-bot/internal/demodata"
	"github.com/euscs/euscs-bot/internal/discord"
	"github.com/euscs/euscs-bot/internal/env"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

func init() {
	err := godotenv.Load("../../.env")
	env.Init()
	if err != nil {
		log.Warning("error loading .env file: " + err.Error())
	}
	log.SetOutput(io.Discard)
	discord.Init()
}

func Test_createNewMatch(t *testing.T) {
	db.Clear()
	db.Init()
	ctx := context.TODO()
	demodata.InitWithBasicInformation()
	t1, _ := db.GetTeamByName(ctx, "team1")
	t3, _ := db.GetTeamByName(ctx, "team3")
	err := createNewMatch(ctx, t1, t3)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
}
