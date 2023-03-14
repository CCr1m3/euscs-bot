package slashcommands

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/euscs/euscs-bot/internal/rank"
	"github.com/euscs/euscs-bot/internal/static"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type Link struct{}

func (p Link) Name() string {
	return "link"
}

func (p Link) Description() string {
	return "Allow you to link to an omega strikers account."
}

func (p Link) RequiredPerm() *int64 {
	perm := int64(discordgo.PermissionSendMessages)
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
	ctx := context.WithValue(context.Background(), static.UUIDKey, uuid.New())
	log.WithFields(log.Fields{
		string(static.UUIDKey):     ctx.Value(static.UUIDKey),
		string(static.CallerIDKey): playerID,
		string(static.UsernameKey): username,
	}).Info("Link slash command invoked")
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Link slash command invoked. Please wait...",
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

	err = rank.LinkPlayerToUsername(ctx, playerID, username)
	if err != nil {
		log.Errorf("failed to link player %s with username %s: "+err.Error(), playerID, username)
		switch {
		case errors.Is(err, static.ErrUsernameInvalid):
			log.WithFields(log.Fields{
				string(static.UUIDKey):     ctx.Value(static.UUIDKey),
				string(static.CallerIDKey): i.Member.User.ID,
				string(static.UsernameKey): username,
			}).Warning("failed to link player, username invalid")
			message = fmt.Sprintf("Failed to link because username does not exist: %s", username)
		case errors.Is(err, static.ErrUserAlreadyLinked):
			log.WithFields(log.Fields{
				string(static.UUIDKey):     ctx.Value(static.UUIDKey),
				string(static.CallerIDKey): i.Member.User.ID,
			}).Warning("failed to link player, user already linked")
			message = "You have already linked an omega strikers account. Please contact a mod if you want to unlink."
		case errors.Is(err, static.ErrUsernameAlreadyLinked):
			log.WithFields(log.Fields{
				string(static.UUIDKey):     ctx.Value(static.UUIDKey),
				string(static.CallerIDKey): i.Member.User.ID,
				string(static.UsernameKey): username,
			}).Warning("failed to link player, username already linked")
			message = fmt.Sprintf("%s is already linked to an account. Please contact a mod if you think you are the rightful owner of the account.", username)
		default:
			log.WithFields(log.Fields{
				string(static.UUIDKey):     ctx.Value(static.UUIDKey),
				string(static.CallerIDKey): i.Member.User.ID,
				string(static.UsernameKey): username,
				string(static.ErrorKey):    err.Error(),
			}).Error("failed to link player")
			message = fmt.Sprintf("Failed to link to %s.", username)
		}
		return
	}
	log.WithFields(log.Fields{
		string(static.UUIDKey):     ctx.Value(static.UUIDKey),
		string(static.CallerIDKey): i.Member.User.ID,
		string(static.UsernameKey): username,
	}).Info("player successfully linked")
	message = fmt.Sprintf("Successfully linked to %s.", username)
}

type Unlink struct{}

func (p Unlink) Name() string {
	return "unlink"
}

func (p Unlink) Description() string {
	return "Allow mods to unlink someone from his omega strikers."
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
	ctx := context.WithValue(context.Background(), static.UUIDKey, uuid.New())
	options := i.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}
	callerID := i.Member.User.ID
	user := optionMap["discorduser"].UserValue(s)
	log.WithFields(log.Fields{
		string(static.UUIDKey):     ctx.Value(static.UUIDKey),
		string(static.CallerIDKey): callerID,
		string(static.PlayerIDKey): user.ID,
	}).Info("unlink slash command invoked")
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Unlink slash command invoked. Please wait...",
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
	if i.Member.Permissions&discordgo.PermissionModerateMembers != discordgo.PermissionModerateMembers {
		message = "You do not have the permission to unlink."
		return
	}
	err = rank.UnlinkPlayer(ctx, user.ID)
	if err != nil {
		if errors.Is(err, static.ErrUserNotLinked) {
			log.WithFields(log.Fields{
				string(static.UUIDKey):     ctx.Value(static.UUIDKey),
				string(static.CallerIDKey): i.Member.User.ID,
				string(static.PlayerIDKey): user.ID,
			}).Warning("player is not linked")
			message = fmt.Sprintf("%s has not linked an omega strikers account.", user.Mention())
		} else {
			log.WithFields(log.Fields{
				string(static.UUIDKey):     ctx.Value(static.UUIDKey),
				string(static.CallerIDKey): i.Member.User.ID,
				string(static.PlayerIDKey): user.ID,
				string(static.ErrorKey):    err.Error(),
			}).Error("failed to unlink player")
			message = fmt.Sprintf("Failed to unlink %s.", user.Mention())
		}
		return
	}
	message = fmt.Sprintf("Successfully unlinked %s.", user.Mention())
}
