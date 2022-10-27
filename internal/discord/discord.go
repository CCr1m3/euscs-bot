package discord

import (
	"os"
	"regexp"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/haashi/omega-strikers-bot/internal/chat"
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
	session.AddHandler(MentionReaction)

	err = session.Open()
	if err != nil {
		log.Fatalf("cannot open the session: %v", err)
	}
	err = initRoles()
	if err != nil {
		log.Fatalf("cannot initialize roles: %v", err)
	}
	go func() {
		threadCleanUp()
		time.Sleep(time.Minute * 5)
	}()
}

func Stop() {
	session.Close()
}

func MentionReaction(s *discordgo.Session, m *discordgo.MessageCreate) {

	r := regexp.MustCompile(s.State.User.Mention())
	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if r.MatchString(m.Content) {
		s.ChannelMessageSend(m.ChannelID, chat.GenerateRandomMessage())
	}
}
