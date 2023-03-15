package discord

import (
	"github.com/bwmarrin/discordgo"
	"github.com/euscs/euscs-bot/internal/env"
	"github.com/euscs/euscs-bot/internal/scheduled"
	log "github.com/sirupsen/logrus"
)

var GuildID string
var session *discordgo.Session

func GetSession() *discordgo.Session {
	return session
}

func Init() {
	log.Info("starting discord service")
	GuildID = env.Discord.GuildID
	botToken := env.Discord.Token
	var err error
	session, err = discordgo.New("Bot " + botToken)
	if err != nil {
		log.Fatalf("invalid bot parameters: %v", err)
	}
	session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})

	err = session.Open()
	if err != nil {
		log.Fatalf("cannot open the session: %v", err)
	}
	err = initRoles()
	if err != nil {
		log.Fatalf("cannot initialize roles: %v", err)
	}
	err = initChannels()
	if err != nil {
		log.Fatalf("cannot initialize channels: %v", err)
	}
}

func Stop() {
	scheduled.TaskManager.Cancel(scheduled.Task{ID: "threadcleanup"})
	session.Close()
}
