package slashcommands

import (
	"errors"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/haashi/omega-strikers-bot/internal/models"
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

func (p Link) RequiredPerm() *int64 {
	perm := int64(discordgo.PermissionViewChannel)
	return &perm
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
		log.Errorf("failed to link player %s with username %s: "+err.Error(), playerID, username)
		var badUsernameErr *models.RankUpdateUsernameError
		var userAlreadyLinkedErr *models.UserAlreadyLinkedError
		var userNameAlreadyLinkedErr *models.UsernameAlreadyLinkedError
		if errors.As(err, &badUsernameErr) {
			message = fmt.Sprintf("Failed to link because username does not exist: %s\n", badUsernameErr.Username)
		} else if errors.As(err, &userAlreadyLinkedErr) {
			message = "You have already linked an omega strikers account. Please contact a mod if you want to unlink."
		} else if errors.As(err, &userNameAlreadyLinkedErr) {
			message = fmt.Sprintf("%s is already linked to an account. Please contact a mod if you think you are the rightful owner of the account.", username)
		} else {
			message = fmt.Sprintf("Failed to link to %s.", username)
		}
	} else {
		message = fmt.Sprintf("Successfully linked to %s.", username)
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: message,
		},
	})
	if err != nil {
		log.Errorf("failed to send message: " + err.Error())
	}
}
