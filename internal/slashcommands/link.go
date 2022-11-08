package slashcommands

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	"github.com/haashi/omega-strikers-bot/internal/models"
	"github.com/haashi/omega-strikers-bot/internal/rank"
	log "github.com/sirupsen/logrus"
)

type Link struct{}

func (p Link) Name() string {
	return "link"
}

func (p Link) Description() string {
	return "Allow you to link to an omega strikers account"
}

func (p Link) RequiredPerm() *int64 {
	perm := int64(discordgo.PermissionViewChannel)
	return &perm
}

func (p Link) Options() []*discordgo.ApplicationCommandOption {
	return []*discordgo.ApplicationCommandOption{
		{
			Name:        "username",
			Description: "User name in omega strikers",
			Type:        discordgo.ApplicationCommandOptionString,
			Required:    true,
		},
	}
}

func (p Link) Run(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}
	playerID := i.Member.User.ID
	username := strings.ToLower(optionMap["username"].StringValue())
	ctx := context.WithValue(context.Background(), models.UUIDKey, uuid.New())
	log.WithFields(log.Fields{
		string(models.UUIDKey):     ctx.Value(models.UUIDKey),
		string(models.CallerIDKey): playerID,
		string(models.UsernameKey): username,
	}).Info("Link slash command invoked")
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Link slash command invoked. Please wait...",
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

	err = rank.LinkPlayerToUsername(ctx, playerID, username)
	if err != nil {
		log.Errorf("failed to link player %s with username %s: "+err.Error(), playerID, username)
		var badUsernameErr *models.RankUpdateUsernameError
		var userAlreadyLinkedErr *models.UserAlreadyLinkedError
		var userNameAlreadyLinkedErr *models.UsernameAlreadyLinkedError
		if errors.As(err, &badUsernameErr) {
			log.WithFields(log.Fields{
				string(models.UUIDKey):     ctx.Value(models.UUIDKey),
				string(models.CallerIDKey): i.Member.User.ID,
				string(models.UsernameKey): username,
			}).Warning("failed to link player, username invalid")
			message = fmt.Sprintf("Failed to link because username does not exist: %s\n", badUsernameErr.Username)
		} else if errors.As(err, &userAlreadyLinkedErr) {
			log.WithFields(log.Fields{
				string(models.UUIDKey):     ctx.Value(models.UUIDKey),
				string(models.CallerIDKey): i.Member.User.ID,
			}).Warning("failed to link player, user already linked")
			message = "You have already linked an omega strikers account. Please contact a mod if you want to unlink."
		} else if errors.As(err, &userNameAlreadyLinkedErr) {
			log.WithFields(log.Fields{
				string(models.UUIDKey):     ctx.Value(models.UUIDKey),
				string(models.CallerIDKey): i.Member.User.ID,
				string(models.UsernameKey): username,
			}).Warning("failed to link player, username already linked")
			message = fmt.Sprintf("%s is already linked to an account. Please contact a mod if you think you are the rightful owner of the account.", username)
		} else {
			log.WithFields(log.Fields{
				string(models.UUIDKey):     ctx.Value(models.UUIDKey),
				string(models.CallerIDKey): i.Member.User.ID,
				string(models.UsernameKey): username,
				string(models.ErrorKey):    err.Error(),
			}).Error("failed to link player")
			message = fmt.Sprintf("Failed to link to %s.", username)
		}
		return
	}
	log.WithFields(log.Fields{
		string(models.UUIDKey):     ctx.Value(models.UUIDKey),
		string(models.CallerIDKey): i.Member.User.ID,
		string(models.UsernameKey): username,
	}).Info("player successfully linked")
	message = fmt.Sprintf("Successfully linked to %s.", username)
}
