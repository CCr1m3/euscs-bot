package slashcommands

import (
	"context"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	"github.com/haashi/omega-strikers-bot/internal/models"
	"github.com/haashi/omega-strikers-bot/internal/team"
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
		},
	}
}
func (p Createteam) Run(s *discordgo.Session, i *discordgo.InteractionCreate) {
	ctx := context.WithValue(context.Background(), models.UUIDKey, uuid.New())
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
		string(models.UUIDKey):     ctx.Value(models.UUIDKey),
		string(models.CallerIDKey): i.Member.User.ID,
		string(models.TeamNameKey): teamName,
	}).Info("createteam slash command invoked")
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Createteam slash command invoked. Please wait...",
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
	if teamName == "" {
		message = "Please enter a teamname."
		log.WithFields(log.Fields{
			string(models.UUIDKey):     ctx.Value(models.UUIDKey),
			string(models.CallerIDKey): i.Member.User.ID,
		}).Warning("createteam failed, no arguments")
		return
	}
	err = team.CreateTeam(ctx, i.Member.User.ID, teamName)
	if err == nil {
		message = fmt.Sprintf("Succesfully created team : %s", teamName)
	} else {
		log.WithFields(log.Fields{
			string(models.UUIDKey):     ctx.Value(models.UUIDKey),
			string(models.CallerIDKey): i.Member.User.ID,
			string(models.ErrorKey):    err.Error(),
			string(models.TeamNameKey): teamName,
		}).Error("failed to created team")
		message = fmt.Sprintf("Failed to created team : %s", teamName)
	}
}
