package slashcommands

import (
	"fmt"
	"net/url"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

type Rank struct{}

func (p Rank) Name() string {
	return "rank"
}

func (p Rank) Description() string {
	return "Allow you to display rank"
}

func (p Rank) Options() []*discordgo.ApplicationCommandOption {
	return []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "username",
			Description: "Omega strikers suername",
			Required:    true,
		},
	}
}

func (p Rank) Run(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("https://corestrike.gg/lookup/%s?region=Europe", url.PathEscape(optionMap["username"].StringValue())),
		},
	})
	if err != nil {
		log.Errorf("failed to send message: " + err.Error())
	}
}
