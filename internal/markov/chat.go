package markov

import (
	"context"
	"regexp"

	"github.com/bwmarrin/discordgo"
	"github.com/euscs/euscs-bot/internal/env"
	"github.com/euscs/euscs-bot/internal/models"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	ctx := context.WithValue(context.Background(), models.UUIDKey, uuid.New())
	if m.Author.ID == s.State.User.ID || env.Discord.GuildID != m.GuildID {
		return
	}
	r := regexp.MustCompile(s.State.User.Mention())
	if r.MatchString(m.Content) {
		messageText := GenerateRandomMessage(ctx)
		rgx := regexp.MustCompile(`<@\d+>`)
		var sanitizedMessageText = rgx.ReplaceAllString(messageText, "@someone")
		mes, err := s.ChannelMessageSend(m.ChannelID, sanitizedMessageText)
		if err != nil {
			log.Error("failed to send message: " + err.Error())
			return
		}
		if messageText != sanitizedMessageText {
			_, err = s.ChannelMessageEdit(m.ChannelID, mes.ID, messageText)
			if err != nil {
				log.Error("failed to edit message: " + err.Error())
				return
			}
		}
	} else {
		err := learn(ctx, m.Content)
		if err != nil {
			log.Error("failed to learn message: " + err.Error())
		}
	}
}
