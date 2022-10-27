package chat

import (
	"os"
	"strings"

	"github.com/haashi/omega-strikers-bot/internal/discord"
	log "github.com/sirupsen/logrus"
)

func fetchAllMessages() {
	session := discord.GetSession()
	f, _ := os.Create("messages")
	channels, err := session.GuildChannels(discord.GuildID)
	if err != nil {
		log.Errorf("failed to get guild channels: " + err.Error())
	}
	for _, channel := range channels {
		//get all messages
		log.Debugf("getting messages from %s", channel.Name)
		messages, err := session.ChannelMessages(channel.ID, 100, "", "", "")
		if err != nil {
			log.Errorf("failed to get messages: " + err.Error())
		}
		for len(messages) != 0 {
			for _, message := range messages {
				if message.Author.ID == session.State.User.ID {
					continue
				}
				_, err = f.WriteString(strings.ToLower(message.Content) + "\n")
				if err != nil {
					log.Errorf("failed to write messages: " + err.Error())
				}
			}
			lastMessage := messages[len(messages)-1]
			messages, err = session.ChannelMessages(channel.ID, 100, lastMessage.ID, "", "")
			if err != nil {
				log.Errorf("failed to get messages: " + err.Error())
			}
		}
	}
}
