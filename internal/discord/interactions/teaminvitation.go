package interactions

import (
	"context"

	"github.com/bwmarrin/discordgo"
	"github.com/euscs/euscs-bot/internal/db"
	"github.com/euscs/euscs-bot/internal/static"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type AcceptInvite struct{}

func (p AcceptInvite) Name() string {
	return "accept_invite"
}

func (p AcceptInvite) Run(s *discordgo.Session, i *discordgo.InteractionCreate) {
	ctx := context.WithValue(context.Background(), static.UUIDKey, uuid.New())
	log.WithFields(log.Fields{
		string(static.UUIDKey):         ctx.Value(static.UUIDKey),
		string(static.CallerIDKey):     i.User.ID,
		string(static.InvitationIDKey): i.Message.ID,
	}).Info("accept invitation interaction invoked")
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Accept invite interaction invoked. Please wait...",
		},
	})
	if err != nil {
		log.WithFields(log.Fields{
			string(static.UUIDKey):  ctx.Value(static.UUIDKey),
			string(static.ErrorKey): err.Error(),
		}).Error("failed to send message")
		return
	}
	invitation, err := db.GetTeamInvitationByID(ctx, i.Message.ID)
	if err != nil {
		log.WithFields(log.Fields{
			string(static.UUIDKey):  ctx.Value(static.UUIDKey),
			string(static.ErrorKey): err.Error(),
		}).Error("failed to get invitation")
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
		s.ChannelMessageEditComplex(&discordgo.MessageEdit{
			Channel:    i.Message.ChannelID,
			ID:         i.Message.ID,
			Components: []discordgo.MessageComponent{},
			Content:    &i.Message.Content,
		})
	}()
	err = invitation.Accept(ctx)
	if err != nil {
		log.WithFields(log.Fields{
			string(static.UUIDKey):  ctx.Value(static.UUIDKey),
			string(static.ErrorKey): err.Error(),
		}).Error("failed to accept invitation")
		message = "Failed to accept invitation."
	} else {
		message = "Successfully accepted invitation."
	}
}

type RefuseInvite struct{}

func (p RefuseInvite) Name() string {
	return "refuse_invite"
}

func (p RefuseInvite) Run(s *discordgo.Session, i *discordgo.InteractionCreate) {
	ctx := context.WithValue(context.Background(), static.UUIDKey, uuid.New())
	log.WithFields(log.Fields{
		string(static.UUIDKey):         ctx.Value(static.UUIDKey),
		string(static.CallerIDKey):     i.User.ID,
		string(static.InvitationIDKey): i.Message.ID,
	}).Info("refuse invitation interaction invoked")
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Refuse invite interaction invoked. Please wait...",
		},
	})
	if err != nil {
		log.WithFields(log.Fields{
			string(static.UUIDKey):  ctx.Value(static.UUIDKey),
			string(static.ErrorKey): err.Error(),
		}).Error("failed to send message")
		return
	}
	invitation, err := db.GetTeamInvitationByID(ctx, i.Message.ID)
	if err != nil {
		log.WithFields(log.Fields{
			string(static.UUIDKey):  ctx.Value(static.UUIDKey),
			string(static.ErrorKey): err.Error(),
		}).Error("failed to get invitation")
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
		s.ChannelMessageEditComplex(&discordgo.MessageEdit{
			Channel:    i.Message.ChannelID,
			ID:         i.Message.ID,
			Components: []discordgo.MessageComponent{},
			Content:    &i.Message.Content,
		})
	}()

	err = invitation.Refuse(ctx)
	if err != nil {
		log.WithFields(log.Fields{
			string(static.UUIDKey):  ctx.Value(static.UUIDKey),
			string(static.ErrorKey): err.Error(),
		}).Error("failed to refused invitation")
		message = "Failed to refused invitation."
	} else {
		message = "Successfully refused invitation."
	}
}
