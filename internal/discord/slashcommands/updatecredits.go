package slashcommands

import (
	"context"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/euscs/euscs-bot/internal/db"
	"github.com/euscs/euscs-bot/internal/static"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type Updatecredits struct{}

func (p Updatecredits) Name() string {
	return "updatecredits"
}

func (p Updatecredits) Description() string {
	return "Allow mods to update one users credits."
}

func (p Updatecredits) RequiredPerm() *int64 {
	perm := int64(discordgo.PermissionModerateMembers)
	return &perm
}

func (p Updatecredits) Options() []*discordgo.ApplicationCommandOption {
	return []*discordgo.ApplicationCommandOption{
		{
			Name:        "discorduser",
			Description: "User in Discord",
			Type:        discordgo.ApplicationCommandOptionUser,
			Required:    true,
		},
		{
			Name:        "option",
			Description: "How to update user's credits",
			Type:        discordgo.ApplicationCommandOptionString,
			Required:    true,
			Choices: []*discordgo.ApplicationCommandOptionChoice{
				{
					Name:  "Add",
					Value: "add",
				},
				{
					Name:  "Subtract",
					Value: "subtract",
				},
				{
					Name:  "Set",
					Value: "set",
				},
			},
		},
		{
			Name:        "amount",
			Description: "How many credits",
			Type:        discordgo.ApplicationCommandOptionInteger,
			Required:    true,
		},
	}
}

func (p Updatecredits) Run(s *discordgo.Session, i *discordgo.InteractionCreate) {
	ctx := context.WithValue(context.Background(), static.UUIDKey, uuid.New())
	options := i.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}
	callerID := i.Member.User.ID
	user := optionMap["discorduser"].UserValue(s)
	option := optionMap["option"].StringValue()
	amount := int(optionMap["amount"].IntValue())
	log.WithFields(log.Fields{
		string(static.UUIDKey):     ctx.Value(static.UUIDKey),
		string(static.CallerIDKey): callerID,
		string(static.PlayerIDKey): user.ID,
	}).Info("updatecredits slash command invoked")
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Updatecredits slash command invoked. Please wait...",
		},
	})
	if err != nil {
		log.WithFields(log.Fields{
			string(static.UUIDKey):  ctx.Value(static.UUIDKey),
			string(static.ErrorKey): err.Error(),
		}).Error("failed to send message")
		return
	}
	var message string
	defer func() {
		_, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: &message,
		})
		if err != nil {
			log.WithFields(log.Fields{
				string(static.UUIDKey):  ctx.Value(static.UUIDKey),
				string(static.ErrorKey): err.Error(),
			}).Error("failed to edit message")
		}
	}()
	if i.Member.Permissions&discordgo.PermissionModerateMembers != discordgo.PermissionModerateMembers {
		message = "You do not have the permission to update credits."
		return
	}
	if amount < 0 {
		message = "Please enter a positive number."
		return
	}

	player, err := db.GetOrCreatePlayerByID(ctx, user.ID)
	if err != nil {
		log.WithFields(log.Fields{
			string(static.UUIDKey):     ctx.Value(static.UUIDKey),
			string(static.ErrorKey):    err.Error(),
			string(static.PlayerIDKey): user.ID,
		}).Error("failed to get or create player")
		return
	}
	switch option {
	case "set":
		player.Credits = amount
		err = player.SetCredits(ctx, player.Credits)
		if err != nil {
			log.Error("failed to update player: " + err.Error())
			message = "Failed to update credits."
			return
		}
	case "add":
		player.Credits += amount
		err = player.SetCredits(ctx, player.Credits)
		if err != nil {
			log.Error("failed to update player: " + err.Error())
			message = "Failed to update credits."
			return
		}
	case "subtract":
		if player.Credits < amount {
			message = "Users can't have negative amount of credits. If you want to remove all credits, use \"set 0\"."
			return
		}
		player.Credits -= amount
		err = player.SetCredits(ctx, player.Credits)
		if err != nil {
			log.Error("failed to update player: " + err.Error())
			message = "Failed to update credits."
			return
		}
	}
	message = fmt.Sprintf("Successfully updated credits of %s.", user.Mention())
	log.Info("updated user: ", user.ID, " ; option: ", option, "; amount: ", amount)
}
