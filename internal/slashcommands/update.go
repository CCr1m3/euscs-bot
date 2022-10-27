package slashcommands

import (
	"errors"

	"github.com/bwmarrin/discordgo"
	"github.com/haashi/omega-strikers-bot/internal/models"
	"github.com/haashi/omega-strikers-bot/internal/rank"
	log "github.com/sirupsen/logrus"
)

type Update struct{}

func (p Update) Name() string {
	return "update"
}

func (p Update) Description() string {
	return "Allow you to update discord role using your linked omega strikers account."
}

func (p Update) RequiredPerm() *int64 {
	perm := int64(discordgo.PermissionViewChannel)
	return &perm
}

func (p Update) Options() []*discordgo.ApplicationCommandOption {
	return []*discordgo.ApplicationCommandOption{}
}

func (p Update) Run(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}
	var message string
	playerID := i.Member.User.ID
	err := rank.UpdateRankIfNeeded(playerID)
	if err != nil {
		log.Errorf("failed to update player %s with username %s: "+err.Error(), playerID)
		var tooFastErr *models.RankUpdateTooFastError
		var notLinkedErr *models.NotLinkedError
		if errors.As(err, &tooFastErr) {
			message = "You have updated your account recently. Please wait before using this command again."
		} else if errors.As(err, &notLinkedErr) {
			message = "You have not linked your omega strikers account. Please use '/rank link' first."
		} else {
			message = "Failed to update your rank."
		}
	} else {
		message = "Successfully updated your rank."
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
