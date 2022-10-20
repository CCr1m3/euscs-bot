package slashcommands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/haashi/omega-strikers-bot/internal/db"
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
					Value: "forward",
				},
				{
					Name:  "Goalie",
					Value: "goalie",
				},
				{
					Name:  "Flex",
					Value: "flex",
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
	player, err := db.GetPlayer(i.Member.User.ID)
	if err != nil {
		log.Infof("failed to get player :" + err.Error())
	}
	if player == nil {
		err = db.CreatePlayer(i.Member.User.ID)
		if err != nil {
			log.Error("failed to create player :" + err.Error())
			return
		}
		player, err = db.GetPlayer(i.Member.User.ID)
		if err != nil {
			log.Error("failed to get player :" + err.Error())
			return
		}
	}
	isInMatch, err := player.IsInMatch()
	if err != nil {
		log.Error("failed to search for player current match :" + err.Error())
		return
	}
	if isInMatch {
		message = "You are already in a match !"
	} else {
		if player.IsInQueue() {
			message = "You are already in the queue !"
		} else {
			err := player.AddToQueue(optionMap["role"].StringValue())
			if err != nil {
				log.Errorf("%s failed to queue :"+err.Error(), player.DiscordID)
				return
			}
			log.Debugf("%s joined the queue as a %s", player.DiscordID, optionMap["role"].StringValue())
			message = fmt.Sprintf("You joined the queue as a %s !", optionMap["role"].StringValue())
		}
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: message,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		log.Error("failed to send message")
	}
}
