package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/euscs/euscs-bot/internal/credits"
	"github.com/euscs/euscs-bot/internal/db"
	"github.com/euscs/euscs-bot/internal/discord"
	"github.com/euscs/euscs-bot/internal/discord/slashcommands"
	"github.com/euscs/euscs-bot/internal/env"
	"github.com/euscs/euscs-bot/internal/markov"
	"github.com/euscs/euscs-bot/internal/matchmaking"
	"github.com/euscs/euscs-bot/internal/random"
	"github.com/euscs/euscs-bot/internal/rank"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

func init() {
	flag.Parse()
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Warning("error loading .env file: " + err.Error())
	}
	env.Init()
	if env.LogLevel == env.DEBUG {
		log.SetLevel(log.DebugLevel)
		log.SetReportCaller(true)
	} else {
		log.SetLevel(log.InfoLevel)
	}
	random.Init()
	db.Init()
	discord.Init()
	markov.Init()
	if env.Mode != env.DEV {
		rank.Init()
	}
	credits.Init()
	slashcommands.Init()
	matchmaking.Init()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	log.Info("initialization done")
	<-stop
	log.Info("gracefully shutting down.")
	discord.Stop()
}
