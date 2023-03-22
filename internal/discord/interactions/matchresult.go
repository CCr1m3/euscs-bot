package interactions

import (
	"context"

	"github.com/bwmarrin/discordgo"
	"github.com/euscs/euscs-bot/internal/static"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type MatchResult struct{}

func (p MatchResult) Name() string {
	return "match_result"
}

func (p MatchResult) Run(s *discordgo.Session, i *discordgo.InteractionCreate) {
	ctx := context.WithValue(context.Background(), static.UUIDKey, uuid.New())
	log.WithFields(log.Fields{
		string(static.UUIDKey):         ctx.Value(static.UUIDKey),
		string(static.CallerIDKey):     i.User.ID,
		string(static.InvitationIDKey): i.Message.ID,
	}).Info("match result interaction invoked")
}
