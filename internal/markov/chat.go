package markov

import (
	"context"
	"math/rand"
	"regexp"

	"github.com/bwmarrin/discordgo"
	"github.com/euscs/euscs-bot/internal/db"
	"github.com/euscs/euscs-bot/internal/env"
	"github.com/euscs/euscs-bot/internal/static"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	ctx := context.WithValue(context.Background(), static.UUIDKey, uuid.New())
	if m.Author.ID == s.State.User.ID || env.Discord.GuildID != m.GuildID {
		return
	}
	r := regexp.MustCompile(s.State.User.Mention())
	if r.MatchString(m.Content) {
		player, err := db.GetOrCreatePlayerByID(ctx, m.Author.ID)
		if err != nil {
			log.Error("failed to get player: " + err.Error())
			return
		}
		if player.Credits >= 10 {
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
			err = player.SetCredits(ctx, player.Credits-10)
			if err != nil {
				log.Error("failed to update player: " + err.Error())
				return
			}
		}
	} else {
		err := learn(ctx, m.Content)
		if err != nil {
			log.Error("failed to learn message: " + err.Error())
		}
		if rand.Intn(10) < 1 {
			player, err := db.GetOrCreatePlayerByID(ctx, m.Author.ID)
			if err != nil {
				log.Error("failed to get player: " + err.Error())
				return
			}
			err = player.SetCredits(ctx, player.Credits+1)
			if err != nil {
				log.Error("failed to update player: " + err.Error())
				return
			}
		}
	}
}
