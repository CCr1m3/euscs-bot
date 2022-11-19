package slashcommands

import (
	"context"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	"github.com/haashi/omega-strikers-bot/internal/models"
	"github.com/haashi/omega-strikers-bot/internal/rank"
	log "github.com/sirupsen/logrus"
)

type Who struct{}

func (p Who) Name() string {
	return "who"
}

func (p Who) Description() string {
	return "Allow you to know about an user"
}

func (p Who) RequiredPerm() *int64 {
	perm := int64(discordgo.PermissionViewChannel)
	return &perm
}

func (p Who) Options() []*discordgo.ApplicationCommandOption {
	return []*discordgo.ApplicationCommandOption{
		{
			Name:        "discorduser",
			Description: "User in discord",
			Type:        discordgo.ApplicationCommandOptionUser,
		},
		{
			Name:        "username",
			Description: "User name in omega strikers",
			Type:        discordgo.ApplicationCommandOptionString,
		},
	}
}

func (p Who) Run(s *discordgo.Session, i *discordgo.InteractionCreate) {
	ctx := context.WithValue(context.Background(), models.UUIDKey, uuid.New())
	options := i.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}
	var userID string
	var username string
	if val, ok := optionMap["discorduser"]; ok {
		userID = val.UserValue(s).ID
	}
	if val, ok := optionMap["username"]; ok {
		username = val.StringValue()
	}
	log.WithFields(log.Fields{
		string(models.UUIDKey):     ctx.Value(models.UUIDKey),
		string(models.CallerIDKey): i.Member.User.ID,
		string(models.UsernameKey): username,
		string(models.PlayerIDKey): userID,
	}).Info("who slash command invoked")
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Who slash command invoked. Please wait...",
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
	if userID == "" && username == "" {
		message = "Please enter at least one of the argument."
		log.WithFields(log.Fields{
			string(models.UUIDKey):     ctx.Value(models.UUIDKey),
			string(models.CallerIDKey): i.Member.User.ID,
		}).Warning("who failed, no arguments")
		return
	}
	if userID != "" && username != "" {
		message = "Please enter only one of the argument."
		log.WithFields(log.Fields{
			string(models.UUIDKey):     ctx.Value(models.UUIDKey),
			string(models.CallerIDKey): i.Member.User.ID,
		}).Warning("who failed, two arguments")
		return
	}
	if userID != "" {
		username, err := rank.GetLinkedUsername(ctx, userID)
		if err != nil {
			log.WithFields(log.Fields{
				string(models.UUIDKey):     ctx.Value(models.UUIDKey),
				string(models.PlayerIDKey): userID,
				string(models.ErrorKey):    err.Error(),
			}).Error("failed to lookup user")
			message = "Failed to lookup user."
			return
		}
		if username == "" {
			log.WithFields(log.Fields{
				string(models.UUIDKey):     ctx.Value(models.UUIDKey),
				string(models.PlayerIDKey): userID,
			}).Warning("player is not linked")
			message = fmt.Sprintf("%s has not linked his omega strikers account.", "<@"+userID+">")
			return
		}
		message = fmt.Sprintf("%s is %s in omega strikers.", "<@"+userID+">", username)
		return
	}
	if username != "" {
		userID, err := rank.GetLinkedUser(ctx, username)
		if err != nil {
			log.WithFields(log.Fields{
				string(models.UUIDKey):     ctx.Value(models.UUIDKey),
				string(models.UsernameKey): username,
				string(models.ErrorKey):    err.Error(),
			}).Error("failed to lookup username")
			message = "Failed to lookup username."
			return
		}
		if userID == "" {
			log.WithFields(log.Fields{
				string(models.UUIDKey):     ctx.Value(models.UUIDKey),
				string(models.UsernameKey): username,
			}).Warning("username is not linked")
			message = fmt.Sprintf("%s is not in this server.", username)
			return
		}
		message = fmt.Sprintf("%s is %s.", username, "<@"+userID+">")
		return
	}
}
