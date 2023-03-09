package slashcommands

import (
	"context"

	"github.com/bwmarrin/discordgo"
	"github.com/euscs/euscs-bot/internal/models"
	"github.com/euscs/euscs-bot/internal/team"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type Invite struct{}

func (p Invite) Name() string {
	return "invite"
}

func (p Invite) Description() string {
	return "Allow you to invite an user"
}

func (p Invite) RequiredPerm() *int64 {
	perm := int64(discordgo.PermissionViewChannel)
	return &perm
}

func (p Invite) Options() []*discordgo.ApplicationCommandOption {
	return []*discordgo.ApplicationCommandOption{
		{
			Name:        "discorduser",
			Description: "User in discord",
			Type:        discordgo.ApplicationCommandOptionUser,
			Required:    true,
		},
	}
}
func (p Invite) Run(s *discordgo.Session, i *discordgo.InteractionCreate) {
	ctx := context.WithValue(context.Background(), models.UUIDKey, uuid.New())
	options := i.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}
	var invitedPlayerID string
	if val, ok := optionMap["discorduser"]; ok {
		invitedPlayerID = val.UserValue(s).ID
	}
	log.WithFields(log.Fields{
		string(models.UUIDKey):     ctx.Value(models.UUIDKey),
		string(models.CallerIDKey): i.Member.User.ID,
		string(models.PlayerIDKey): invitedPlayerID,
	}).Info("invite slash command invoked")
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Invite slash command invoked. Please wait...",
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
	if invitedPlayerID == "" {
		message = "Please enter user in discord."
		log.WithFields(log.Fields{
			string(models.UUIDKey):     ctx.Value(models.UUIDKey),
			string(models.CallerIDKey): i.Member.User.ID,
		}).Warning("invite failed, no arguments")
		return
	}
	err = team.InvitePlayer(ctx, i.Member.User.ID, invitedPlayerID)
	if err == nil {
		message = "Successfully invited."
	} else {
		log.WithFields(log.Fields{
			string(models.UUIDKey):     ctx.Value(models.UUIDKey),
			string(models.CallerIDKey): i.Member.User.ID,
			string(models.PlayerIDKey): invitedPlayerID,
		}).Warning("invite failed, no arguments")
		message = "Failed to invite."
	}
}
