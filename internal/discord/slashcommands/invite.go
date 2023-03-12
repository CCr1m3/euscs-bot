package slashcommands

import (
	"context"
	"errors"

	"github.com/bwmarrin/discordgo"
	"github.com/euscs/euscs-bot/internal/static"
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
	ctx := context.WithValue(context.Background(), static.UUIDKey, uuid.New())
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
		string(static.UUIDKey):     ctx.Value(static.UUIDKey),
		string(static.CallerIDKey): i.Member.User.ID,
		string(static.PlayerIDKey): invitedPlayerID,
	}).Info("invite slash command invoked")
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Invite slash command invoked. Please wait...",
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
	if invitedPlayerID == "" {
		message = "Please enter user in discord."
		log.WithFields(log.Fields{
			string(static.UUIDKey):     ctx.Value(static.UUIDKey),
			string(static.CallerIDKey): i.Member.User.ID,
		}).Warning("invite failed, no arguments")
		return
	}
	err = team.InvitePlayerToTeam(ctx, i.Member.User.ID, invitedPlayerID)
	if err == nil {
		message = "Successfully invited."
	} else {
		message = "Failed to invite."
		log.WithFields(log.Fields{
			string(static.UUIDKey):     ctx.Value(static.UUIDKey),
			string(static.CallerIDKey): i.Member.User.ID,
			string(static.PlayerIDKey): invitedPlayerID,
			string(static.ErrorKey):    err.Error(),
		}).Warning("invite failed")
		switch {
		case errors.Is(err, static.ErrTeamFull):
			message += " Your team is full."
		case errors.Is(err, static.ErrUserAlreadyInTeam):
			message += " This user already has a team."
		case errors.Is(err, static.ErrNotFound):
			message += " You don't have a team."
		case errors.Is(err, static.ErrNotTeamOwner):
			message += " You are not the team owner."
		default:
			message += "Unexpected Error."
		}
	}
}
