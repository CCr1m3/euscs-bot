package discord

import (
	"os"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

var GuildID string
var BotToken string
var session *discordgo.Session

func GetSession() *discordgo.Session {
	return session
}

func Init() {
	log.Info("starting discord service")
	GuildID = os.Getenv("guildid")
	BotToken = os.Getenv("token")
	var err error
	session, err = discordgo.New("Bot " + BotToken)
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

}

func Stop() {
	session.Close()
}
