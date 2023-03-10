package team

import (
	"context"
	"errors"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/euscs/euscs-bot/internal/db"
	"github.com/euscs/euscs-bot/internal/discord"
	"github.com/euscs/euscs-bot/internal/models"
)

func CreateTeam(ctx context.Context, ownerID string, teamName string) error {
	player, err := db.GetOrCreatePlayerById(ctx, ownerID)
	if err != nil {
		return err
	}
	team := models.Team{Players: []*models.Player{player}, OwnerID: ownerID, Name: teamName}
	err = db.CreateTeam(ctx, &team)
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
		return models.ErrTeamFull
	}
	session := discord.GetSession()
	channel, err := session.UserChannelCreate(playerID)
	if err != nil {
		return err
	}
	team2, err := db.GetTeamByPlayerID(ctx, playerID)
	if err != nil {
		if !errors.Is(err, models.ErrNotFound) {
			return err
		}
	}
	if team2 != nil {
		return models.ErrUserAlreadyInTeam
	}

	message := fmt.Sprintf("You have been invited by %s to the team '%s'", "<@"+ownerID+">", team.Name)
	_, err = session.ChannelMessageSendComplex(channel.ID, &discordgo.MessageSend{
		Content: message,
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
	return nil
}

func AddPlayerToTeam(ctx context.Context, teamName string, playerID string) error {
	team, err := db.GetTeamByName(ctx, teamName)
	if err != nil {
		return err
	}
	if len(team.Players) >= 3 {
		return models.ErrTeamFull
	}
	return nil
}
