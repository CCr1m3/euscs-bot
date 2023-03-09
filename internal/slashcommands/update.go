package slashcommands

import (
	"context"
	"errors"

	"github.com/bwmarrin/discordgo"
	"github.com/euscs/euscs-bot/internal/models"
	"github.com/euscs/euscs-bot/internal/rank"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type Update struct{}

func (p Update) Name() string {
	return "update"
}

func (p Update) Description() string {
	return "Allow you to update discord role using your linked omega strikers account."
}

func (p Update) RequiredPerm() *int64 {
	perm := int64(discordgo.PermissionViewChannel)
	return &perm
}

func (p Update) Options() []*discordgo.ApplicationCommandOption {
	return []*discordgo.ApplicationCommandOption{}
}

func (p Update) Run(s *discordgo.Session, i *discordgo.InteractionCreate) {
	ctx := context.WithValue(context.Background(), models.UUIDKey, uuid.New())
	playerID := i.Member.User.ID
	log.WithFields(log.Fields{
		string(models.UUIDKey):     ctx.Value(models.UUIDKey),
		string(models.CallerIDKey): i.Member.User.ID,
	}).Info("update slash command invoked")
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Update slash command invoked. Please wait...",
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

	err = rank.UpdateRankIfNeeded(ctx, playerID)
	if err != nil {
		var tooFastErr *models.RankUpdateTooFastError
		var notLinkedErr *models.NotLinkedError
		if errors.As(err, &tooFastErr) {
			log.WithFields(log.Fields{
				string(models.UUIDKey):     ctx.Value(models.UUIDKey),
				string(models.CallerIDKey): i.Member.User.ID,
			}).Warning("player update too fast")
			message = "You have updated your account recently. Please wait before using this command again."
		} else if errors.As(err, &notLinkedErr) {
			log.WithFields(log.Fields{
				string(models.UUIDKey):     ctx.Value(models.UUIDKey),
				string(models.CallerIDKey): i.Member.User.ID,
			}).Warning("player is not linked")
			message = "You have not linked your omega strikers account. Please use 'link' first."
		} else {
			log.WithFields(log.Fields{
				string(models.UUIDKey):     ctx.Value(models.UUIDKey),
				string(models.CallerIDKey): i.Member.User.ID,
				string(models.ErrorKey):    err.Error(),
			}).Error("failed to update player")
			message = "Failed to update your rank."
		}
		return
	}
	message = "Successfully updated your rank."
}
