package slashcommands

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	"github.com/haashi/omega-strikers-bot/internal/credits"
	"github.com/haashi/omega-strikers-bot/internal/matchmaking"
	"github.com/haashi/omega-strikers-bot/internal/models"
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
	ctx := context.WithValue(context.Background(), models.UUIDKey, uuid.New())
	log.WithFields(log.Fields{
		string(models.UUIDKey):      ctx.Value(models.UUIDKey),
		string(models.CallerIDKey):  i.Member.User.ID,
		string(models.ChannelIDKey): i.ChannelID,
		string(models.ResultKey):    team,
	}).Info("predict slash command invoked")
	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "predict slash command invoked. Please wait...",
		},
	})
	if err != nil {
		log.WithFields(log.Fields{
			string(models.UUIDKey):  ctx.Value(models.UUIDKey),
			string(models.ErrorKey): err.Error(),
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
				string(models.UUIDKey):  ctx.Value(models.UUIDKey),
				string(models.ErrorKey): err.Error(),
			}).Error("failed to edit message")
		}
	}()

	match, err := matchmaking.GetMatchByThreadId(ctx, i.ChannelID)
	if err != nil {
		log.WithFields(log.Fields{
			string(models.UUIDKey):      ctx.Value(models.UUIDKey),
			string(models.ChannelIDKey): i.ChannelID,
			string(models.ErrorKey):     err.Error(),
		}).Warning("failed to find match")
		message = "This channel is not a match lobby."
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
		log.WithFields(log.Fields{
			string(models.UUIDKey):     ctx.Value(models.UUIDKey),
			string(models.CallerIDKey): i.Member.User.ID,
			string(models.MatchIDKey):  match.ID,
		}).Warning("can't predict, user is in match")
		message = "You are a player of this match."
		return
	}
	if time.Since(time.Unix(int64(match.Timestamp), 0)) > time.Minute*3 {
		log.WithFields(log.Fields{
			string(models.UUIDKey):     ctx.Value(models.UUIDKey),
			string(models.CallerIDKey): i.Member.User.ID,
			string(models.MatchIDKey):  match.ID,
			string(models.DurationKey): time.Since(time.Unix(int64(match.Timestamp), 0)),
		}).Warning("can't predict, not in time")
		message = "The match has already started for too long to predict."
		return
	}
	err = credits.AddPrediction(ctx, i.Member.User.ID, match.ID, int(team))
	if err != nil {
		log.WithFields(log.Fields{
			string(models.UUIDKey):     ctx.Value(models.UUIDKey),
			string(models.CallerIDKey): i.Member.User.ID,
			string(models.MatchIDKey):  match.ID,
		}).Warning("can't predict, already predicted")
		message = "You have already predicted in this match."
		return
	}
	message = fmt.Sprintf("%s predicted team%d victory.", i.Member.User.Mention(), team)
}
