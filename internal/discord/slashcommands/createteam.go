package slashcommands

import (
	"context"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/euscs/euscs-bot/internal/db"
	"github.com/euscs/euscs-bot/internal/static"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type Createteam struct{}

func (p Createteam) Name() string {
	return "createteam"
}

func (p Createteam) Description() string {
	return "Allow you to create a team"
}

func (p Createteam) RequiredPerm() *int64 {
	perm := int64(discordgo.PermissionViewChannel)
	return &perm
}

func (p Createteam) Options() []*discordgo.ApplicationCommandOption {
	return []*discordgo.ApplicationCommandOption{
		{
			Name:        "teamname",
			Description: "name of the team",
			Type:        discordgo.ApplicationCommandOptionString,
			Required:    true,
		},
	}
}
func (p Createteam) Run(s *discordgo.Session, i *discordgo.InteractionCreate) {
	ctx := context.WithValue(context.Background(), static.UUIDKey, uuid.New())
	options := i.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}
	var teamName string
	if val, ok := optionMap["teamname"]; ok {
		teamName = val.StringValue()
	}
	log.WithFields(log.Fields{
		string(static.UUIDKey):     ctx.Value(static.UUIDKey),
		string(static.CallerIDKey): i.Member.User.ID,
		string(static.TeamNameKey): teamName,
	}).Info("createteam slash command invoked")
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Createteam slash command invoked. Please wait...",
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
	if teamName == "" {
		message = "Please enter a teamname."
		log.WithFields(log.Fields{
			string(static.UUIDKey):     ctx.Value(static.UUIDKey),
			string(static.CallerIDKey): i.Member.User.ID,
		}).Warning("createteam failed, no arguments")
		return
	}
	player, err := db.GetOrCreatePlayerByID(ctx, i.Member.User.ID)
	if err != nil {
		log.WithFields(log.Fields{
			string(static.UUIDKey):     ctx.Value(static.UUIDKey),
			string(static.CallerIDKey): i.Member.User.ID,
			string(static.ErrorKey):    err.Error(),
		}).Error("failed to get player")
		message = fmt.Sprintf("Failed to created team : %s", teamName)
	}
	_, err = player.CreateTeamWithName(ctx, teamName)
	if err == nil {
		message = fmt.Sprintf("Succesfully created team : %s", teamName)
	} else {
		log.WithFields(log.Fields{
			string(static.UUIDKey):     ctx.Value(static.UUIDKey),
			string(static.CallerIDKey): i.Member.User.ID,
			string(static.ErrorKey):    err.Error(),
			string(static.TeamNameKey): teamName,
		}).Error("failed to created team")
		message = fmt.Sprintf("Failed to created team : %s", teamName)
	}
}
