package slashcommands

import (
	"context"

	"github.com/bwmarrin/discordgo"
	"github.com/euscs/euscs-bot/internal/db"
	"github.com/euscs/euscs-bot/internal/matchmaking"
	"github.com/euscs/euscs-bot/internal/static"
	"github.com/google/uuid"
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
	ctx := context.WithValue(context.Background(), static.UUIDKey, uuid.New())
	log.WithFields(log.Fields{
		string(static.UUIDKey):      ctx.Value(static.UUIDKey),
		string(static.CallerIDKey):  i.Member.User.ID,
		string(static.ChannelIDKey): i.ChannelID,
	}).Info("cancel slash command invoked")
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Cancel slash command invoked. Please wait...",
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
			string(static.UUIDKey):     ctx.Value(static.UUIDKey),
			string(static.CallerIDKey): i.Member.User.ID,
			string(static.MatchIDKey):  match.ID,
		}).Warning("can't cancel, user not in match")
		message = "You are not a player of this match."
		return
	}
	if match.State == db.MatchStateVoteInProgress {
		log.WithFields(log.Fields{
			string(static.UUIDKey):    ctx.Value(static.UUIDKey),
			string(static.MatchIDKey): match.ID,
		}).Warning("can't cancel, match vote already in progress")
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
	message = "Confirmation vote started."
	matchmaking.VoteCancelMatch(ctx, match)
}
