package slashcommands

import (
	"context"
	"errors"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/euscs/euscs-bot/internal/db"
	"github.com/euscs/euscs-bot/internal/matchmaking"
	"github.com/euscs/euscs-bot/internal/static"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type Join struct{}

func (p Join) Name() string {
	return "join"
}

func (p Join) Description() string {
	return "Allow you to join the queue"
}

func (p Join) RequiredPerm() *int64 {
	perm := int64(discordgo.PermissionViewChannel)
	return &perm
}

func (p Join) Options() []*discordgo.ApplicationCommandOption {
	return []*discordgo.ApplicationCommandOption{
		{
			Name:        "role",
			Description: "Role in omega strikers",
			Type:        discordgo.ApplicationCommandOptionString,
			Required:    true,
			Choices: []*discordgo.ApplicationCommandOptionChoice{
				{
					Name:  "Forward",
					Value: db.RoleForward,
				},
				{
					Name:  "Goalie",
					Value: db.RoleGoalie,
				},
				{
					Name:  "Flex",
					Value: db.RoleFlex,
				},
			},
		},
	}
}

func (p Join) Run(s *discordgo.Session, i *discordgo.InteractionCreate) {
	ctx := context.WithValue(context.Background(), static.UUIDKey, uuid.New())
	options := i.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}
	playerID := i.Member.User.ID
	role := optionMap["role"].StringValue()
	log.WithFields(log.Fields{
		string(static.UUIDKey):      ctx.Value(static.UUIDKey),
		string(static.CallerIDKey):  i.Member.User.ID,
		string(static.QueueRoleKey): role,
	}).Info("join slash command invoked")
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Join slash command invoked. Please wait...",
			Flags:   discordgo.MessageFlagsEphemeral,
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

	isInQueue, err := matchmaking.IsPlayerInQueue(ctx, playerID)
	if err != nil {
		log.WithFields(log.Fields{
			string(static.UUIDKey):     ctx.Value(static.UUIDKey),
			string(static.CallerIDKey): i.Member.User.ID,
			string(static.ErrorKey):    err.Error(),
		}).Error("failed to check if player is in queue")
		message = "Failed to put you in queue."
		return
	}
	isInMatch, err := matchmaking.IsPlayerInMatch(ctx, playerID)
	if err != nil {
		log.WithFields(log.Fields{
			string(static.UUIDKey):     ctx.Value(static.UUIDKey),
			string(static.CallerIDKey): i.Member.User.ID,
			string(static.ErrorKey):    err.Error(),
		}).Error("failed to check if player is in match")
		message = "Failed to put you in queue."
		return
	}
	if isInMatch {
		message = "You are already in a match !"
		return
	}
	if isInQueue {
		message = "You are already in the queue !"
		return
	}
	err = matchmaking.AddPlayerToQueue(ctx, playerID, db.Role(role))
	if err != nil {
		if errors.Is(err, static.ErrUserNotLinked) {
			log.WithFields(log.Fields{
				string(static.UUIDKey):     ctx.Value(static.UUIDKey),
				string(static.CallerIDKey): i.Member.User.ID,
			}).Warning("player is not yet linked")
			message = "You have not linked your Omega Strikers account. Please use '/link' first."
		} else {
			log.WithFields(log.Fields{
				string(static.UUIDKey):     ctx.Value(static.UUIDKey),
				string(static.CallerIDKey): i.Member.User.ID,
				string(static.ErrorKey):    err.Error(),
			}).Error("failed to put player in the queue")
			message = "Failed to put you in the queue."
		}
		return
	}
	log.WithFields(log.Fields{
		string(static.UUIDKey):      ctx.Value(static.UUIDKey),
		string(static.CallerIDKey):  i.Member.User.ID,
		string(static.QueueRoleKey): role,
	}).Info("player joined the queue")
	message = fmt.Sprintf("You joined the queue as a %s!", role)
}
