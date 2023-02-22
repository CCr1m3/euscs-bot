package main

import (
	"flag"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/haashi/omega-strikers-bot/internal/credits"
	"github.com/haashi/omega-strikers-bot/internal/db"
	"github.com/haashi/omega-strikers-bot/internal/discord"
	"github.com/haashi/omega-strikers-bot/internal/markov"
	"github.com/haashi/omega-strikers-bot/internal/matchmaking"
	"github.com/haashi/omega-strikers-bot/internal/slashcommands"
	"github.com/haashi/omega-strikers-bot/internal/webserver"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

func init() {
	flag.Parse()
}

func main() {
	rand.Seed(time.Now().Unix())
	err := godotenv.Load(".env")
	if err != nil {
		log.Warning("error loading .env file: " + err.Error())
	}
	logLevel := os.Getenv("loglevel")
	if logLevel == "debug" {
		log.SetLevel(log.DebugLevel)
		log.SetReportCaller(true)
	} else {
		log.SetLevel(log.InfoLevel)
	}
	db.Init()
	discord.Init()
	markov.Init()
	credits.Init()
	//slashcommands.Init()
	matchmaking.Init()
	webserver.Init()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	log.Info("initialization done")
	<-stop
	log.Info("gracefully shutting down.")
	slashcommands.Stop()
	discord.Stop()
}
