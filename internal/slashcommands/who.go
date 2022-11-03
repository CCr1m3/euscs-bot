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
	user := optionMap["discorduser"].UserValue(s)
	username := optionMap["username"].StringValue()
	log.WithFields(log.Fields{
		string(models.UUIDKey):     ctx.Value(models.UUIDKey),
		string(models.CallerIDKey): i.Member.User.ID,
		string(models.UsernameKey): user,
		string(models.PlayerIDKey): user.ID,
	}).Info("who slash command invoked")
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Who slash command invoked. Please wait...",
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

	if user != nil && username != "" {
		message = "Please enter only one of the argument."
		log.WithFields(log.Fields{
			string(models.UUIDKey):     ctx.Value(models.UUIDKey),
			string(models.CallerIDKey): i.Member.User.ID,
		}).Warning("who failed, no arguments")
		return
	}
	if user != nil {
		username, err := rank.GetLinkedUsername(user.ID)
		if err != nil {
			log.WithFields(log.Fields{
				string(models.UUIDKey):     ctx.Value(models.UUIDKey),
				string(models.PlayerIDKey): user.ID,
				string(models.ErrorKey):    err.Error(),
			}).Error("failed to lookup user")
			message = "Failed to lookup user."
			return
		}
		if username == "" {
			log.WithFields(log.Fields{
				string(models.UUIDKey):     ctx.Value(models.UUIDKey),
				string(models.PlayerIDKey): user.ID,
			}).Warning("player is not linked")
			message = fmt.Sprintf("%s has not linked his omega strikers account.", "<@"+user.ID+">")
			return
		}
		message = fmt.Sprintf("%s is %s in omega strikers.", "<@"+user.ID+">", username)
		return
	}
	if username != "" {
		userID, err := rank.GetLinkedUser(username)
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
