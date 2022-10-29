package slashcommands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/haashi/omega-strikers-bot/internal/currency"
	log "github.com/sirupsen/logrus"
)

type Currency struct{}

func (p Currency) Name() string {
	return "currency"
}

func (p Currency) Description() string {
	return "currency description"
}

func (p Currency) RequiredPerm() *int64 {
	perm := int64(discordgo.PermissionViewChannel)
	return &perm
}

func (p Currency) Options() []*discordgo.ApplicationCommandOption {
	return nil
}

func (p Currency) Run(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var message string
	playerID := i.Member.User.ID
	currency, err := currency.GetPlayerCurrency(playerID)
	if err != nil {
		log.Errorf("failed to get player %s currency: "+err.Error(), playerID)
		message = "Failed to get your currency."
	} else {
		message = fmt.Sprintf("You have : %d", currency)
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
