package match

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/euscs/euscs-bot/internal/db"
	"github.com/euscs/euscs-bot/internal/discord"
	"github.com/google/uuid"
)

func createNewMatch(ctx context.Context, teamA *db.Team, teamB *db.Team) error {
	matchId := uuid.New().String()
	channelId := discord.MatchesChannel.ID
	session := discord.GetSession()
	match := &db.Match{}
	match.ID = matchId
	match.Timestamp = int(time.Now().Unix())
	mentionMsg := ""
	teamAPlayersMessage := ""
	for i := range teamA.Players {
		teamAPlayersMessage += "<@" + teamA.Players[i].DiscordID + ">\n"
		mentionMsg += "<@" + teamA.Players[i].DiscordID + ">"
	}
	teamBPlayersMessage := ""
	for i := range teamB.Players {
		teamBPlayersMessage += "<@" + teamB.Players[i].DiscordID + ">\n"
		mentionMsg += "<@" + teamB.Players[i].DiscordID + ">"
	}
	embed := discordgo.MessageEmbed{
		Title: "| Match " + matchId,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "A: " + teamA.Name,
				Value:  teamAPlayersMessage,
				Inline: true,
			},
			{
				Name:   "B: " + teamB.Name,
				Value:  teamBPlayersMessage,
				Inline: true,
			},
			{
				Name:  "State",
				Value: "In progress",
			},
		},
	}
	initialMessage, err := session.ChannelMessageSendComplex(channelId, &discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{
			&embed,
		},
	})
	if err != nil {
		return err
	}
	match.MessageID = initialMessage.ID
	thread, err := session.MessageThreadStartComplex(initialMessage.ChannelID, initialMessage.ID, &discordgo.ThreadStart{
		Name:                matchId,
		AutoArchiveDuration: 1440,
		Invitable:           true,
	})
	if err != nil {
		return err
	}
	match.ThreadID = thread.ID
	lobbyCode := strings.Split(matchId, "-")[0]
	_, err = session.ChannelMessageSendComplex(thread.ID, &discordgo.MessageSend{
		Content: fmt.Sprintf("%s\nLobby code : %s\n\nUse this thread to chat.\nPlease report match result by clicking on the WINNER team button.", mentionMsg, lobbyCode),
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						Label:    teamA.Name,
						Style:    discordgo.PrimaryButton,
						CustomID: "match_result_A",
					},
					discordgo.Button{
						Label:    teamB.Name,
						Style:    discordgo.PrimaryButton,
						CustomID: "match_result_B",
					},
				},
			},
		},
	})
	if err != nil {
		return err
	}
	match.Team1 = teamA
	match.Team2 = teamB
	return nil
}
