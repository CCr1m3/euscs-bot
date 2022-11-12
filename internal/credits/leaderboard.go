package credits

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/haashi/omega-strikers-bot/internal/db"
	"github.com/haashi/omega-strikers-bot/internal/discord"
	"github.com/haashi/omega-strikers-bot/internal/models"
	"github.com/haashi/omega-strikers-bot/internal/scheduled"
	log "github.com/sirupsen/logrus"
)

func Init() {
	scheduled.TaskManager.Add(scheduled.Task{ID: "updateLeaderboard", Run: updateLeaderboard, Frequency: time.Minute})
}

func updateLeaderboard() {
	ctx := context.WithValue(context.Background(), models.UUIDKey, uuid.New())
	players, err := db.GetPlayersOrderedByCredits(ctx)
	if err != nil {
		log.WithFields(log.Fields{
			string(models.UUIDKey):  ctx.Value(models.UUIDKey),
			string(models.ErrorKey): err.Error(),
		}).Error("failed to get players ordered by credits")
		return
	}
	nbMessagesNeeded := len(players)/10 + 1
	session := discord.GetSession()
	messages, err := session.ChannelMessages(discord.LeaderboardChannel.ID, 100, "", "", "")
	if err != nil {
		log.WithFields(log.Fields{
			string(models.UUIDKey):  ctx.Value(models.UUIDKey),
			string(models.ErrorKey): err.Error(),
		}).Error("failed to get messages from leaderboard channel")
		return
	}
	for i := 0; i < (nbMessagesNeeded - len(messages)); i++ {
		_, err := session.ChannelMessageSend(discord.LeaderboardChannel.ID, "placeholder")
		if err != nil {
			log.WithFields(log.Fields{
				string(models.UUIDKey):  ctx.Value(models.UUIDKey),
				string(models.ErrorKey): err.Error(),
			}).Error("failed to send messages in leaderboard channel")
			return
		}
	}
	messages, err = session.ChannelMessages(discord.LeaderboardChannel.ID, 100, "", "", "")
	if err != nil {
		log.WithFields(log.Fields{
			string(models.UUIDKey):  ctx.Value(models.UUIDKey),
			string(models.ErrorKey): err.Error(),
		}).Error("failed to get messages from leaderboard channel")
		return
	}
	contents := make([]string, nbMessagesNeeded)
	for i, player := range players {
		contents[i/10] += fmt.Sprintf("%d: %s (%d credits)\n", i, "<@"+player.DiscordID+">", player.Credits)
	}
	for i := 0; i < nbMessagesNeeded; i++ {
		_, err := session.ChannelMessageEdit(discord.LeaderboardChannel.ID, messages[i].ID, contents[nbMessagesNeeded-i-1])
		if err != nil {
			log.WithFields(log.Fields{
				string(models.UUIDKey):  ctx.Value(models.UUIDKey),
				string(models.ErrorKey): err.Error(),
			}).Error("failed to edit messages from leaderboard channel")
			return
		}
	}
}
