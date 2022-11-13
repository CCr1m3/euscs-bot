package markov

import (
	"context"
	"math/rand"
	"regexp"

	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	"github.com/haashi/omega-strikers-bot/internal/db"
	"github.com/haashi/omega-strikers-bot/internal/discord"
	"github.com/haashi/omega-strikers-bot/internal/models"
	log "github.com/sirupsen/logrus"
)

func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	ctx := context.WithValue(context.Background(), models.UUIDKey, uuid.New())
	if m.Author.ID == s.State.User.ID || discord.GuildID != m.GuildID {
		return
	}
	r := regexp.MustCompile(s.State.User.Mention())
	if r.MatchString(m.Content) {
		player, err := db.GetOrCreatePlayerById(ctx, m.Author.ID)
		if err != nil {
			log.Error("failed to get player: " + err.Error())
			return
		}
		if player.Credits >= 10 {
			player.Credits -= 10
			_, err := s.ChannelMessageSend(m.ChannelID, GenerateRandomMessage(ctx))
			if err != nil {
				log.Error("failed to send message: " + err.Error())
				return
			}
			err = db.UpdatePlayer(ctx, player)
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
			player, err := db.GetOrCreatePlayerById(ctx, m.Author.ID)
			if err != nil {
				log.Error("failed to get player: " + err.Error())
				return
			}
			player.Credits += 1
			err = db.UpdatePlayer(ctx, player)
			if err != nil {
				log.Error("failed to update player: " + err.Error())
				return
			}
		}
	}
}
