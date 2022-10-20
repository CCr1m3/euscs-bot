package slashcommands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/haashi/omega-strikers-bot/internal/db"
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
	player, err := db.GetPlayer(i.Member.User.ID)
	if err != nil {
		log.Infof("failed to get player :" + err.Error())
	}
	if player == nil {
		err = db.CreatePlayer(i.Member.User.ID)
		if err != nil {
			log.Errorf("failed to create player :" + err.Error())
			return
		}
		player, err = db.GetPlayer(i.Member.User.ID)
		if err != nil {
			log.Errorf("failed to get player :" + err.Error())
			return
		}
	}
	if !player.IsInQueue() {
		message = "You are not in the queue !"
	} else {
		err = player.LeaveQueue()
		if err != nil {
			log.Errorf("%s failed to leave the queue :"+err.Error(), player.DiscordID)
			return
		}
		log.Debugf("%s left the queue", player.DiscordID)
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
		log.Error("failed to send message")
	}
}
