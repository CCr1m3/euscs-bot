package slashcommands

import (
	"fmt"
	"math"

	"github.com/bwmarrin/discordgo"
	"github.com/haashi/omega-strikers-bot/internal/matchmaking"
	"github.com/haashi/omega-strikers-bot/internal/models"
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
	perm := int64(discordgo.PermissionViewChannel)
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

	log.Debugf("%s used /result on channel %s with parameters: (%d-%d)", i.Member.User.ID, i.ChannelID, team1Score, team2Score)
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
	if match.State == models.MatchStateVoteInProgress {
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "A confirmation is already in progress.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if err != nil {
			log.Error("failed to send message: " + err.Error())
		}
		return
	}
	if math.Abs(float64(team1Score-team2Score)) < 2 {
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("The result (%d-%d) is not a valid result.", team1Score, team2Score),
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
			Content: "Confirmation started.",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		log.Error("failed to send message: " + err.Error())
	}
	matchmaking.VoteResultMatch(match, int(team1Score), int(team2Score))
}
