package main

import (
	"log"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func main() {
	guildid := "1023624205272617061"
	bottoken := ""
	var err error
	session, err := discordgo.New(bottoken)
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
	f, _ := os.Create("log")
	channels, err := session.GuildChannels(guildid)
	for _, channel := range channels {
		//get all messages
		messages, err := session.ChannelMessages(channel.ID, 100, "", "", "")
		if err != nil {
			log.Fatalf("invalid bot parameters: %v", err)
		}
		for len(messages) != 0 {
			for _, message := range messages {
				f.WriteString(strings.ToLower(message.Content))
				f.WriteString("\n")
			}
			lastMessage := messages[len(messages)-1]
			messages, err = session.ChannelMessages(channel.ID, 100, lastMessage.ID, "", "")
			if err != nil {
				log.Fatalf("invalid bot parameters: %v", err)
			}
		}
	}
}
