package slashcommands

import (
	"context"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/euscs/euscs-bot/internal/credits"
	"github.com/euscs/euscs-bot/internal/static"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type Credits struct{}

func (p Credits) Name() string {
	return "credits"
}

func (p Credits) Description() string {
	return "Allow you to know how many credits you have."
}

func (p Credits) RequiredPerm() *int64 {
	perm := int64(discordgo.PermissionViewChannel)
	return &perm
}

func (p Credits) Options() []*discordgo.ApplicationCommandOption {
	return nil
}

func (p Credits) Run(s *discordgo.Session, i *discordgo.InteractionCreate) {
	playerID := i.Member.User.ID
	ctx := context.WithValue(context.Background(), static.UUIDKey, uuid.New())
	log.WithFields(log.Fields{
		string(static.UUIDKey):     ctx.Value(static.UUIDKey),
		string(static.CallerIDKey): playerID,
	}).Info("credits slash command invoked")
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Credits slash command invoked. Please wait...",
		},
	})
	if err != nil {
		log.WithFields(log.Fields{
			string(static.UUIDKey):  ctx.Value(static.UUIDKey),
			string(static.ErrorKey): err.Error(),
		}).Error("failed to send message")
		return
	}
	var message string
	defer func() {
		_, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: &message,
		})
		if err != nil {
			log.WithFields(log.Fields{
				string(static.UUIDKey):  ctx.Value(static.UUIDKey),
				string(static.ErrorKey): err.Error(),
			}).Error("failed to edit message")
		}
	}()

	credits, err := credits.GetPlayerCredits(ctx, playerID)
	if err != nil {
		log.WithFields(log.Fields{
			string(static.UUIDKey):     ctx.Value(static.UUIDKey),
			string(static.CallerIDKey): i.Member.User.ID,
			string(static.ErrorKey):    err.Error(),
		}).Error("failed to get player credits")
		message = "Failed to get your credits."
		return
	}
	message = fmt.Sprintf("You have %d Credits.", credits)
}
