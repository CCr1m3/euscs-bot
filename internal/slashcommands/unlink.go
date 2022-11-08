package slashcommands

import (
	"context"
	"errors"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	"github.com/haashi/omega-strikers-bot/internal/models"
	"github.com/haashi/omega-strikers-bot/internal/rank"
	log "github.com/sirupsen/logrus"
)

type Unlink struct{}

func (p Unlink) Name() string {
	return "unlink"
}

func (p Unlink) Description() string {
	return "Allow mods to unlink someone from his omega strikers"
}

func (p Unlink) RequiredPerm() *int64 {
	perm := int64(discordgo.PermissionModerateMembers)
	return &perm
}

func (p Unlink) Options() []*discordgo.ApplicationCommandOption {
	return []*discordgo.ApplicationCommandOption{
		{
			Name:        "discorduser",
			Description: "User in discord",
			Type:        discordgo.ApplicationCommandOptionUser,
			Required:    true,
		},
	}
}

func (p Unlink) Run(s *discordgo.Session, i *discordgo.InteractionCreate) {
	ctx := context.WithValue(context.Background(), models.UUIDKey, uuid.New())
	options := i.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}
	callerID := i.Member.User.ID
	user := optionMap["discorduser"].UserValue(s)
	log.WithFields(log.Fields{
		string(models.UUIDKey):     ctx.Value(models.UUIDKey),
		string(models.CallerIDKey): callerID,
		string(models.PlayerIDKey): user.ID,
	}).Info("unlink slash command invoked")
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Unlink slash command invoked. Please wait...",
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

	if i.Member.Permissions&discordgo.PermissionModerateMembers != discordgo.PermissionModerateMembers {
		message = "You do not have the permission to unlink."
		return
	}
	err = rank.UnlinkPlayer(ctx, user.ID)
	if err != nil {
		var notLinkedErr *models.NotLinkedError
		if errors.As(err, &notLinkedErr) {
			log.WithFields(log.Fields{
				string(models.UUIDKey):     ctx.Value(models.UUIDKey),
				string(models.CallerIDKey): i.Member.User.ID,
				string(models.PlayerIDKey): user.ID,
			}).Warning("player is not linked")
			message = fmt.Sprintf("%s has not linked an omega strikers account.", user.Mention())
		} else {
			log.WithFields(log.Fields{
				string(models.UUIDKey):     ctx.Value(models.UUIDKey),
				string(models.CallerIDKey): i.Member.User.ID,
				string(models.PlayerIDKey): user.ID,
				string(models.ErrorKey):    err.Error(),
			}).Error("failed to unlink player")
			message = fmt.Sprintf("Failed to unlink %s.", user.Mention())
		}
		return
	}
	message = fmt.Sprintf("Successfully unlink %s.", user.Mention())
}
