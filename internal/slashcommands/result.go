package slashcommands

import (
	"fmt"
	"math"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/haashi/omega-strikers-bot/internal/matchmaking"
	log "github.com/sirupsen/logrus"
)

type Result struct{}

func (p Result) Name() string {
	return "result"
}

func (p Result) Description() string {
	return "Allow you to report a result using scores : team1 vs team2"
}

func (p Result) RequiredPerm() *int64 {
	perm := int64(discordgo.PermissionSendMessages)
	return &perm
}

func (p Result) Options() []*discordgo.ApplicationCommandOption {
	return []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionInteger,
			Name:        "team1-score",
			Description: "Score",
			Required:    true,
		},
		{
			Type:        discordgo.ApplicationCommandOptionInteger,
			Name:        "team2-score",
			Description: "Score",
			Required:    true,
		},
	}
}

func (p Result) Run(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}
	team1Score := optionMap["team1-score"].IntValue()
	team2Score := optionMap["team2-score"].IntValue()

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
	if math.Abs(float64(team1Score-team2Score)) < 2 {
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("The result (%d-%d) is not a valid result.", team1Score, team2Score),
			},
		})
		if err != nil {
			log.Fatal("failed to send message")
		}
		return
	}
	var message string = fmt.Sprintf("User %s reported score : (%d-%d).\nPlease react to this message to confirm score.", i.Member.Mention(), optionMap["team1-score"].IntValue(), optionMap["team2-score"].IntValue())
	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: message,
		},
	})
	if err != nil {
		log.Error("failed to send message")
	}
	discMessage, err := s.InteractionResponse(i.Interaction)
	if err != nil {
		log.Error("failed to get message")
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
			err = matchmaking.CloseMatch(match, int(team1Score), int(team2Score))
			if err != nil {
				log.Errorf("failed to close match %s: "+err.Error(), match.ID)
			}
			return
		} else if len(playersNOK) > requiredReactions {
			err = s.ChannelMessageDelete(discMessage.ChannelID, discMessage.ID)
			if err != nil {
				log.Errorf("failed to delete message: " + err.Error())
			}
			return
		}
		time.Sleep(time.Second)
	}
}
