package slashcommands

import (
	"context"
	"fmt"
	"math"

	"github.com/bwmarrin/discordgo"
	"github.com/euscs/euscs-bot/internal/db"
	"github.com/euscs/euscs-bot/internal/matchmaking"
	"github.com/euscs/euscs-bot/internal/static"
	"github.com/google/uuid"
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
	ctx := context.WithValue(context.Background(), static.UUIDKey, uuid.New())
	log.WithFields(log.Fields{
		string(static.UUIDKey):      ctx.Value(static.UUIDKey),
		string(static.CallerIDKey):  i.Member.User.ID,
		string(static.ChannelIDKey): i.ChannelID,
		string(static.ResultKey):    fmt.Sprintf("%d-%d", team1Score, team2Score),
	}).Info("result slash command invoked")
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Result slash command invoked. Please wait...",
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

	match, err := matchmaking.GetMatchByThreadId(ctx, i.ChannelID)
	if err != nil {
		log.WithFields(log.Fields{
			string(static.UUIDKey):      ctx.Value(static.UUIDKey),
			string(static.ChannelIDKey): i.ChannelID,
			string(static.ErrorKey):     err.Error(),
		}).Warning("failed to find match")
		message = "This channel is not a match lobby."
		return
	}
	if (team1Score > 3 || team1Score < 0) || (team2Score > 3 || team2Score < 0) {
		log.WithFields(log.Fields{
			string(static.UUIDKey):      ctx.Value(static.UUIDKey),
			string(static.ChannelIDKey): i.ChannelID,
		}).Warning("incorrect result values")
		message = "Please enter proper result values ranging from 0 to 3."
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
	if !inMatch {
		log.WithFields(log.Fields{
			string(static.UUIDKey):      ctx.Value(static.UUIDKey),
			string(static.ChannelIDKey): i.ChannelID,
		}).Warning("can't result, user not in match")
		message = "You are not a player of this match."
		return
	}
	if match.State == db.MatchStateVoteInProgress {
		log.WithFields(log.Fields{
			string(static.UUIDKey):      ctx.Value(static.UUIDKey),
			string(static.ChannelIDKey): i.ChannelID,
		}).Warning("can't result, confirmation already in progress")
		message = "A confirmation is already in progress."
		return
	}
	if match.State != db.MatchStateInProgress {
		log.WithFields(log.Fields{
			string(static.UUIDKey):      ctx.Value(static.UUIDKey),
			string(static.ChannelIDKey): i.ChannelID,
		}).Warning("can't result, match is over")
		message = "The match is already over."
		return
	}
	if math.Abs(float64(team1Score-team2Score)) < 2 {
		message = fmt.Sprintf("The result (%d-%d) is not a valid result.", team1Score, team2Score)
		return
	}
	message = "Confirmation started."
	matchmaking.VoteResultMatch(ctx, match, int(team1Score), int(team2Score))
}
