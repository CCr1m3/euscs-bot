package slashcommands

import (
	"errors"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/haashi/omega-strikers-bot/internal/models"
	"github.com/haashi/omega-strikers-bot/internal/rank"
	log "github.com/sirupsen/logrus"
)

type Unlink struct{}

func (p Unlink) Name() string {
	return "unlink"
}

func (p Unlink) Description() string {
	return "Allow mods to unlink to an omega strikers account"
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
	options := i.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}
	var message string

	if i.Member.Permissions&discordgo.PermissionModerateMembers != discordgo.PermissionModerateMembers {
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "You do not have the permission to unlink.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if err != nil {
			log.Errorf("failed to send message:" + err.Error())
		}
	} else {
		user := optionMap["discorduser"].UserValue(s)
		playerID := user.ID

		err := rank.UnlinkPlayer(playerID)
		if err != nil {
			log.Errorf("failed to unlink player %s: "+err.Error(), playerID)
			var notLinkedErr *models.NotLinkedError
			if errors.As(err, &notLinkedErr) {
				message = fmt.Sprintf("%s has not linked an omega strikers account.", user.Mention())
			} else {
				message = fmt.Sprintf("Failed to unlink %s.", user.Mention())
			}
		} else {
			message = fmt.Sprintf("Successfully unlink %s.", user.Mention())
		}

		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: message,
			},
		})
		if err != nil {
			log.Errorf("failed to send message: " + err.Error())
		}
	}
}
