package slashcommands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/haashi/omega-strikers-bot/internal/matchmaking"
	"github.com/haashi/omega-strikers-bot/internal/models"
	log "github.com/sirupsen/logrus"
)

type Join struct{}

func (p Join) Name() string {
	return "join"
}

func (p Join) Description() string {
	return "Allow you to join the queue"
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
					Value: models.RoleForward,
				},
				{
					Name:  "Goalie",
					Value: models.RoleGoalie,
				},
				{
					Name:  "Flex",
					Value: models.RoleFlex,
				},
			},
		},
	}
}

func (p Join) Run(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}
	var message string
	playerID := i.Member.User.ID

	isInQueue, err := matchmaking.IsPlayerInQueue(playerID)
	if err != nil {
		log.Errorf("failed to check if player is in queue: " + err.Error())
	}
	isInMatch, err := matchmaking.IsPlayerInMatch(playerID)
	if err != nil {
		log.Errorf("failed to check if player is in match: " + err.Error())
	}
	if isInMatch {
		message = "You are already in a match !"
	} else if isInQueue {
		message = "You are already in the queue !"
	} else {
		err = matchmaking.AddPlayerToQueue(playerID, models.Role(optionMap["role"].StringValue()))
		if err != nil {
			log.Errorf("%s failed to queue: "+err.Error(), playerID)
			return
		}
		message = fmt.Sprintf("You joined the queue as a %s !", optionMap["role"].StringValue())
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
