package slashcommands

import (
	"context"

	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	"github.com/haashi/omega-strikers-bot/internal/matchmaking"
	"github.com/haashi/omega-strikers-bot/internal/models"
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
	ctx := context.WithValue(context.Background(), models.UUIDKey, uuid.New())
	log.WithFields(log.Fields{
		string(models.UUIDKey):      ctx.Value(models.UUIDKey),
		string(models.CallerIDKey):  i.Member.User.ID,
		string(models.ChannelIDKey): i.ChannelID,
	}).Info("cancel slash command invoked")
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Cancel slash command invoked. Please wait...",
			Flags:   discordgo.MessageFlagsEphemeral,
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
	if !inMatch {
		log.WithFields(log.Fields{
			string(models.UUIDKey):     ctx.Value(models.UUIDKey),
			string(models.CallerIDKey): i.Member.User.ID,
			string(models.MatchIDKey):  match.ID,
		}).Warning("can't cancel, user not in match")
		message = "You are not a player of this match."
		return
	}
	if match.State == models.MatchStateVoteInProgress {
		log.WithFields(log.Fields{
			string(models.UUIDKey):    ctx.Value(models.UUIDKey),
			string(models.MatchIDKey): match.ID,
		}).Warning("can't cancel, match vote already in progress")
		message = "A confirmation is already in progress."
		return
	}
	if match.State != models.MatchStateInProgress {
		log.WithFields(log.Fields{
			string(models.UUIDKey):      ctx.Value(models.UUIDKey),
			string(models.ChannelIDKey): i.ChannelID,
			string(models.ErrorKey):     err.Error(),
		}).Warning("can't result, match is over")
		message = "The match is already over."
		return
	}
	message = "Confirmation vote started."
	matchmaking.VoteCancelMatch(ctx, match)
}
