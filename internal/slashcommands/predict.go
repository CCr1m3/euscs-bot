package slashcommands

import (
	"fmt"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/haashi/omega-strikers-bot/internal/currency"
	"github.com/haashi/omega-strikers-bot/internal/matchmaking"
	log "github.com/sirupsen/logrus"
)

type Predict struct{}

func (p Predict) Name() string {
	return "predict"
}

func (p Predict) Description() string {
	return "Allow you to predict on a match."
}

func (p Predict) RequiredPerm() *int64 {
	perm := int64(discordgo.PermissionViewChannel)
	return &perm
}

func (p Predict) Options() []*discordgo.ApplicationCommandOption {
	return []*discordgo.ApplicationCommandOption{
		{
			Name:        "team",
			Description: "Which team will win",
			Type:        discordgo.ApplicationCommandOptionString,
			Required:    true,
			Choices: []*discordgo.ApplicationCommandOptionChoice{
				{
					Name:  "Team1",
					Value: "1",
				},
				{
					Name:  "Team2",
					Value: "2",
				},
			},
		},
	}
}

func (p Predict) Run(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}
	team, err := strconv.ParseInt(optionMap["team"].StringValue(), 10, 0)
	if err != nil {
		return
	}
	log.Debugf("%s used /predict on channel with parameters: %d", i.Member.User.ID, team)
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
			log.Error("failed to send message: " + err.Error())
		}
		return
	}
	inMatch := false
	for _, p := range match.Team1 {
		if p.DiscordID == i.Member.User.ID {
			inMatch = true
		}
	}
	for _, p := range match.Team2 {
		if p.DiscordID == i.Member.User.ID {
			inMatch = true
		}
	}
	if inMatch {
		log.Warningf("can't predict: user %s is in match %s", i.Member.User.ID, match.ID)
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "You are a player of this match.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if err != nil {
			log.Error("failed to send message: " + err.Error())
		}
		return
	}
	if time.Since(time.Unix(int64(match.Timestamp), 0)) > time.Minute {
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "The match has already started for too long to predict.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if err != nil {
			log.Error("failed to send message: " + err.Error())
		}
		return
	}
	err = currency.AddPrediction(i.Member.User.ID, match.ID, int(team))
	if err != nil {
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "You already predicted for this match.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if err != nil {
			log.Error("failed to send message: " + err.Error())
		}
		return
	}
	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("%s predicted team%d victory.", i.Member.User.Mention(), team),
		},
	})
	if err != nil {
		log.Error("failed to send message: " + err.Error())
	}
}
