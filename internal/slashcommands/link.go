package slashcommands

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/haashi/omega-strikers-bot/internal/rank"
	log "github.com/sirupsen/logrus"
)

type Link struct{}

func (p Link) Name() string {
	return "link"
}

func (p Link) Description() string {
	return "Allow you to link to an omega strikers account"
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
	var message string
	playerID := i.Member.User.ID
	username := strings.ToLower(optionMap["username"].StringValue())
	err := rank.LinkPlayerToUsername(playerID, username)
	if err != nil {
		log.Errorf("failed to link player %s with username %s :"+err.Error(), playerID, username)
		message = fmt.Sprintf("Failed to link to %s.", username)
	} else {
		message = fmt.Sprintf("Successfully linked to %s.", username)
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: message,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		log.Errorf("failed to send message:" + err.Error())
	}
}
