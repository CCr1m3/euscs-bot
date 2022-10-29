package slashcommands

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/haashi/omega-strikers-bot/internal/rank"
	log "github.com/sirupsen/logrus"
)

type Who struct{}

func (p Who) Name() string {
	return "who"
}

func (p Who) Description() string {
	return "Allow you to know about an user"
}

func (p Who) RequiredPerm() *int64 {
	perm := int64(discordgo.PermissionViewChannel)
	return &perm
}

func (p Who) Options() []*discordgo.ApplicationCommandOption {
	return []*discordgo.ApplicationCommandOption{
		{
			Name:        "discorduser",
			Description: "User in discord",
			Type:        discordgo.ApplicationCommandOptionUser,
		},
		{
			Name:        "username",
			Description: "User name in omega strikers",
			Type:        discordgo.ApplicationCommandOptionString,
		},
	}
}

func (p Who) Run(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}
	var username string
	var user *discordgo.User
	if option, ok := optionMap["username"]; ok {
		username = option.StringValue()
		username = strings.ToLower(username)
	}
	if option, ok := optionMap["discorduser"]; ok {
		user = option.UserValue(s)
	}

	var message string

	if user != nil && username != "" {
		message = "Please enter only one of the argument."
	} else if user != nil {

		username, err := rank.GetLinkedUsername(user.ID)
		if err != nil {
			log.Errorf("failed to get username of %s: "+err.Error(), user)
		}
		if username == "" {
			message = fmt.Sprintf("%s has not linked his omega strikers account.", "<@"+user.ID+">")
		} else {

			message = fmt.Sprintf("%s is %s in omega strikers.", "<@"+user.ID+">", username)
		}
	} else if username != "" {
		userID, err := rank.GetLinkedUser(username)
		if err != nil {
			log.Errorf("failed to get user of %s: "+err.Error(), user)
		}
		if userID == "" {
			message = fmt.Sprintf("%s is not in this server.", username)
		} else {
			message = fmt.Sprintf("%s is %s.", username, "<@"+userID+">")
		}
	} else {
		message = "Please enter at least one of the argument."
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: message,
		},
	})
	if err != nil {
		log.Errorf("failed to send message:" + err.Error())
	}
}
