package slashcommands

import (
	"fmt"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/haashi/omega-strikers-bot/internal/matchmaking"
	log "github.com/sirupsen/logrus"
)

type Cancel struct{}

func (p Cancel) Name() string {
	return "cancel"
}

func (p Cancel) Description() string {
	return "Allow you to cancel a match."
}

func (p Cancel) RequiredPerm() *int64 {
	perm := int64(discordgo.PermissionViewChannel)
	return &perm
}

func (p Cancel) Options() []*discordgo.ApplicationCommandOption {
	return []*discordgo.ApplicationCommandOption{}
}

func (p Cancel) Run(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	log.Debugf("%s used /cancel on channel %s", i.Member.User.ID, i.ChannelID)
	match, err := matchmaking.GetMatchByThreadId(i.ChannelID)
	if err != nil {
		log.Warningf("failed to find match by threadID %s : "+err.Error(), i.ChannelID)
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "This channel is not a match lobby.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if err != nil {
			log.Fatal("failed to send message")
		}
		return
	}
	var message string = fmt.Sprintf("User %s wants to cancel this match.\nPlease react to this message to confirm.", i.Member.Mention())
	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: message,
		},
	})
	log.Debugf("getting players confirmation of cancellation of match %s", match.ID)
	if err != nil {
		log.Error("failed to send message: " + err.Error())
	}
	discMessage, err := s.InteractionResponse(i.Interaction)
	if err != nil {
		log.Error("failed to get message: " + err.Error())
	}
	err = s.MessageReactionAdd(discMessage.ChannelID, discMessage.ID, "✅")
	if err != nil {
		log.Errorf("failed add reaction: " + err.Error())
	}
	err = s.MessageReactionAdd(discMessage.ChannelID, discMessage.ID, "❌")
	if err != nil {
		log.Errorf("failed add reaction: " + err.Error())
	}
	for {
		playersOK, err := s.MessageReactions(discMessage.ChannelID, discMessage.ID, "✅", 10, "", "")
		if err != nil {
			log.Errorf("failed to get reactions: " + err.Error())
		}
		playersNOK, err := s.MessageReactions(discMessage.ChannelID, discMessage.ID, "❌", 10, "", "")
		if err != nil {
			log.Errorf("failed to get reactions: " + err.Error())
		}
		requiredReactions := 4
		if os.Getenv("mode") == "dev" {
			requiredReactions = 1
		}
		if len(playersOK) > requiredReactions {
			log.Debugf("players confirmed cancellation of match %s", match.ID)
			err = matchmaking.CancelMatch(match)
			if err != nil {
				log.Errorf("failed to cancel match %s: "+err.Error(), match.ID)
			}
			return
		} else if len(playersNOK) > requiredReactions {
			log.Debugf("players refused cancellation of match %s", match.ID)
			err = s.ChannelMessageDelete(discMessage.ChannelID, discMessage.ID)
			if err != nil {
				log.Errorf("failed to delete message: " + err.Error())
			}
			return
		}
		time.Sleep(time.Second)
	}
}
