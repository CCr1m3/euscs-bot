package slashcommands

import (
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

type Ping struct{}

func (p Ping) Name() string {
	return "ping"
}

func (p Ping) Description() string {
	return "ping description"
}

func (p Ping) RequiredPerm() *int64 {
	perm := int64(discordgo.PermissionSendMessages)
	return &perm
}

func (p Ping) Options() []*discordgo.ApplicationCommandOption {
	return nil
}

func (p Ping) Run(s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Pong",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		log.Errorf("failed to send message: " + err.Error())
	}
}
