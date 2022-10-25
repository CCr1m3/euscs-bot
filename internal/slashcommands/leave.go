package slashcommands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/haashi/omega-strikers-bot/internal/matchmaking"
	log "github.com/sirupsen/logrus"
)

type Leave struct{}

func (p Leave) Name() string {
	return "leave"
}

func (p Leave) Description() string {
	return "Allow you to leave the queue"
}

func (p Leave) Options() []*discordgo.ApplicationCommandOption {
	return nil
}

func (p Leave) Run(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var message string
	playerID := i.Member.User.ID

	isInQueue, err := matchmaking.IsPlayerInQueue(playerID)
	if err != nil {
		log.Errorf("failed to check if player is in queue: " + err.Error())
	}
	if !isInQueue {
		message = "You are not in the queue !"
	} else {
		err = matchmaking.RemovePlayerFromQueue(playerID)
		if err != nil {
			log.Errorf("%s failed to leave the queue: "+err.Error(), playerID)
			return
		}
		message = "You left the queue !"
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
