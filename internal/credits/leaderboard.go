package credits

import (
	"context"
	"fmt"
	"time"

	"github.com/euscs/euscs-bot/internal/db"
	"github.com/euscs/euscs-bot/internal/discord"
	"github.com/euscs/euscs-bot/internal/scheduled"
	"github.com/euscs/euscs-bot/internal/static"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

func updateLeaderboard() {
	ctx := context.WithValue(context.Background(), static.UUIDKey, uuid.New())
	players, err := db.GetPlayersOrderedByCredits(ctx)
	if err != nil {
		log.WithFields(log.Fields{
			string(static.UUIDKey):  ctx.Value(static.UUIDKey),
			string(static.ErrorKey): err.Error(),
		}).Error("failed to get players ordered by credits")
		return
	}
	nbMessagesNeeded := (len(players)-1)/7 + 1
	maxMessages := 100
	session := discord.GetSession()
	messages, err := session.ChannelMessages(discord.LeaderboardChannel.ID, maxMessages, "", "", "")
	if err != nil {
		log.WithFields(log.Fields{
			string(static.UUIDKey):  ctx.Value(static.UUIDKey),
			string(static.ErrorKey): err.Error(),
		}).Error("failed to get messages from leaderboard channel")
		return
	}
	for i := 0; i < (nbMessagesNeeded - len(messages)); i++ {
		_, err := session.ChannelMessageSend(discord.LeaderboardChannel.ID, "placeholder")
		if err != nil {
			log.WithFields(log.Fields{
				string(static.UUIDKey):  ctx.Value(static.UUIDKey),
				string(static.ErrorKey): err.Error(),
			}).Error("failed to send messages in leaderboard channel")
			return
		}
	}
	for i := 0; i < (len(messages) - nbMessagesNeeded); i++ {
		err := session.ChannelMessageDelete(discord.LeaderboardChannel.ID, messages[i].ID)
		if err != nil {
			log.WithFields(log.Fields{
				string(static.UUIDKey):  ctx.Value(static.UUIDKey),
				string(static.ErrorKey): err.Error(),
			}).Error("failed to delete messages in leaderboard channel")
			return
		}
	}
	messages, err = session.ChannelMessages(discord.LeaderboardChannel.ID, 100, "", "", "")
	if err != nil {
		log.WithFields(log.Fields{
			string(static.UUIDKey):  ctx.Value(static.UUIDKey),
			string(static.ErrorKey): err.Error(),
		}).Error("failed to get messages from leaderboard channel")
		return
	}
	contents := make([]string, nbMessagesNeeded)
	for i, player := range players {
		contents[i/5] += fmt.Sprintf("%d: %s (%d credits)\n", i+1, "<@"+player.DiscordID+">", player.Credits)
	}
	for i := 0; i < nbMessagesNeeded; i++ {
		_, err := session.ChannelMessageEdit(discord.LeaderboardChannel.ID, messages[i].ID, contents[nbMessagesNeeded-i-1])
		if err != nil {
			log.WithFields(log.Fields{
				string(static.UUIDKey):  ctx.Value(static.UUIDKey),
				string(static.ErrorKey): err.Error(),
			}).Error("failed to edit messages from leaderboard channel")
			return
		}
	}
}

func Init() {
	scheduled.TaskManager.Add(scheduled.Task{ID: "updateLeaderboard", Run: updateLeaderboard, Frequency: time.Minute * 5})
}
