package slashcommands

import (
	"context"

	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	"github.com/haashi/omega-strikers-bot/internal/matchmaking"
	"github.com/haashi/omega-strikers-bot/internal/models"
	log "github.com/sirupsen/logrus"
)

type Leave struct{}

func (p Leave) Name() string {
	return "leave"
}

func (p Leave) Description() string {
	return "Allow you to leave the queue"
}

func (p Leave) RequiredPerm() *int64 {
	perm := int64(discordgo.PermissionViewChannel)
	return &perm
}

func (p Leave) Options() []*discordgo.ApplicationCommandOption {
	return nil
}

func (p Leave) Run(s *discordgo.Session, i *discordgo.InteractionCreate) {
	ctx := context.WithValue(context.Background(), models.UUIDKey, uuid.New())
	playerID := i.Member.User.ID
	log.WithFields(log.Fields{
		string(models.UUIDKey):     ctx.Value(models.UUIDKey),
		string(models.CallerIDKey): i.Member.User.ID,
	}).Info("leave slash command invoked")
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Leave slash command invoked. Please wait...",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		log.WithFields(log.Fields{
			string(models.UUIDKey):  ctx.Value(models.UUIDKey),
			string(models.ErrorKey): err.Error(),
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
				string(models.UUIDKey):  ctx.Value(models.UUIDKey),
				string(models.ErrorKey): err.Error(),
			}).Error("failed to edit message")
		}
	}()

	isInQueue, err := matchmaking.IsPlayerInQueue(ctx, playerID)
	if err != nil {
		log.WithFields(log.Fields{
			string(models.UUIDKey):     ctx.Value(models.UUIDKey),
			string(models.CallerIDKey): i.Member.User.ID,
			string(models.ErrorKey):    err.Error(),
		}).Error("failed to check if player is in queue")
		message = "Failed to make you leave the queue."
		return
	}
	if !isInQueue {
		message = "You are not in the queue !"
		return
	}
	err = matchmaking.RemovePlayerFromQueue(ctx, playerID)
	if err != nil {
		log.WithFields(log.Fields{
			string(models.UUIDKey):     ctx.Value(models.UUIDKey),
			string(models.CallerIDKey): i.Member.User.ID,
			string(models.ErrorKey):    err.Error(),
		}).Error("failed to check if player is in queue")
		message = "Failed to make you leave the queue."
		return
	}
	message = "You left the queue !"
}
