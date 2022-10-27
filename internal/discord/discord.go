package discord

import (
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/haashi/omega-strikers-bot/internal/scheduled"
	log "github.com/sirupsen/logrus"
)

var GuildID string
var session *discordgo.Session

func GetSession() *discordgo.Session {
	return session
}

func Init() {
	log.Info("starting discord service")
	GuildID = os.Getenv("guildid")
	botToken := os.Getenv("token")
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
	scheduled.TaskManager.Add(scheduled.Task{ID: "threadcleanup", Run: threadCleanUp, Frequency: time.Hour})
}

func Stop() {
	scheduled.TaskManager.Cancel(scheduled.Task{ID: "threadcleanup"})
	session.Close()
}
