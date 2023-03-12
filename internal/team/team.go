package team

import (
	"context"
	"errors"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/euscs/euscs-bot/internal/db"
	"github.com/euscs/euscs-bot/internal/discord"
	"github.com/euscs/euscs-bot/internal/static"
)

func CreateTeam(ctx context.Context, ownerID string, teamName string) error {
	player, err := db.GetOrCreatePlayerByID(ctx, ownerID)
	if err != nil {
		return err
	}
	team := db.Team{Players: []*db.Player{player}, OwnerID: ownerID, Name: teamName}
	err = team.Save(ctx)
	if err != nil {
		return err
	}
	return nil
}

func InvitePlayerToTeam(ctx context.Context, ownerID string, playerID string) error {
	team, err := db.GetTeamByPlayerID(ctx, ownerID)
	if err != nil {
		return fmt.Errorf("team: %w", err)
	}
	if len(team.Players) >= 3 {
		return static.ErrTeamFull
	}
	if team.OwnerID != ownerID {
		return static.ErrNotTeamOwner
	}
	player, err := db.GetOrCreatePlayerByID(ctx, playerID)
	if err != nil {
		return err
	}
	team2, err := db.GetTeamByPlayerID(ctx, playerID)
	if err != nil {
		if !errors.Is(err, static.ErrNotFound) {
			return err
		}
	}
	if team2 != nil {
		return static.ErrUserAlreadyInTeam
	}
	session := discord.GetSession()
	channel, err := session.UserChannelCreate(playerID)
	if err != nil {
		return err
	}
	messageContent := fmt.Sprintf("You have been invited by %s to the team '%s'", "<@"+ownerID+">", team.Name)
	message, err := session.ChannelMessageSendComplex(channel.ID, &discordgo.MessageSend{
		Content: messageContent,
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						Label:    "Accept",
						Style:    discordgo.PrimaryButton,
						CustomID: "accept_invite",
					},
					discordgo.Button{
						Label:    "Refuse",
						Style:    discordgo.DangerButton,
						CustomID: "refuse_invite",
					},
				},
			},
		}})
	if err != nil {
		return err
	}
	invite := db.TeamInvitation{Player: player, Team: team, MessageID: message.ID, Timestamp: int(message.Timestamp.Unix()), State: db.InvitationPending}
	err = invite.Save(ctx)
	if err != nil {
		return err
	}
	return nil
}

func AcceptInvitation(ctx context.Context, teamName string, playerID string) error {
	team, err := db.GetTeamByName(ctx, teamName)
	if err != nil {
		return err
	}
	if len(team.Players) >= 3 {
		return static.ErrTeamFull
	}
	return nil
}
