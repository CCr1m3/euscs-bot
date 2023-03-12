package slashcommands

import (
	"context"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/euscs/euscs-bot/internal/rank"
	"github.com/euscs/euscs-bot/internal/static"
	"github.com/google/uuid"
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
	ctx := context.WithValue(context.Background(), static.UUIDKey, uuid.New())
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
		string(static.UUIDKey):     ctx.Value(static.UUIDKey),
		string(static.CallerIDKey): i.Member.User.ID,
		string(static.UsernameKey): username,
		string(static.PlayerIDKey): userID,
	}).Info("who slash command invoked")
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Who slash command invoked. Please wait...",
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
	if userID == "" && username == "" {
		message = "Please enter at least one of the argument."
		log.WithFields(log.Fields{
			string(static.UUIDKey):     ctx.Value(static.UUIDKey),
			string(static.CallerIDKey): i.Member.User.ID,
		}).Warning("who failed, no arguments")
		return
	}
	if userID != "" && username != "" {
		message = "Please enter only one of the argument."
		log.WithFields(log.Fields{
			string(static.UUIDKey):     ctx.Value(static.UUIDKey),
			string(static.CallerIDKey): i.Member.User.ID,
		}).Warning("who failed, two arguments")
		return
	}
	if userID != "" {
		username, err := rank.GetLinkedUsername(ctx, userID)
		if err != nil {
			log.WithFields(log.Fields{
				string(static.UUIDKey):     ctx.Value(static.UUIDKey),
				string(static.PlayerIDKey): userID,
				string(static.ErrorKey):    err.Error(),
			}).Error("failed to lookup user")
			message = "Failed to lookup user."
			return
		}
		if username == "" {
			log.WithFields(log.Fields{
				string(static.UUIDKey):     ctx.Value(static.UUIDKey),
				string(static.PlayerIDKey): userID,
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
				string(static.UUIDKey):     ctx.Value(static.UUIDKey),
				string(static.UsernameKey): username,
				string(static.ErrorKey):    err.Error(),
			}).Error("failed to lookup username")
			message = "Failed to lookup username."
			return
		}
		if userID == "" {
			log.WithFields(log.Fields{
				string(static.UUIDKey):     ctx.Value(static.UUIDKey),
				string(static.UsernameKey): username,
			}).Warning("username is not linked")
			message = fmt.Sprintf("%s is not in this server.", username)
			return
		}
		message = fmt.Sprintf("%s is %s.", username, "<@"+userID+">")
		return
	}
}
