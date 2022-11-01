package slashcommands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/haashi/omega-strikers-bot/internal/credits"
	log "github.com/sirupsen/logrus"
)

type Credits struct{}

func (p Credits) Name() string {
	return "credits"
}

func (p Credits) Description() string {
	return "Allow you to know how many credits you have."
}

func (p Credits) RequiredPerm() *int64 {
	perm := int64(discordgo.PermissionViewChannel)
	return &perm
}

func (p Credits) Options() []*discordgo.ApplicationCommandOption {
	return nil
}

func (p Credits) Run(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var message string
	playerID := i.Member.User.ID
	credits, err := credits.GetPlayerCredits(playerID)
	if err != nil {
		log.Errorf("failed to get player %s credits: "+err.Error(), playerID)
		message = "Failed to get your credits."
	} else {
		message = fmt.Sprintf("You have : %d", credits)
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: message,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		log.Errorf("failed to send message: " + err.Error())
	}
}
