package credits

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
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
	nbMessagesNeeded := (len(players)-1)/5 + 1
	maxMessages := 100

	// count current msgs in #credits-leaderboard
	lenMsgs := 0
	session := discord.GetSession()
	messages, err := session.ChannelMessages(discord.LeaderboardChannel.ID, maxMessages, "", "", "")
	if err != nil {
		if strings.Contains(err.Error(), "500 Internal Server Error") {
			log.Info("failed to request messages" + err.Error())
			log.Info("reattempting request...")
			time.Sleep(time.Second / 10)
			messages, err = session.ChannelMessages(discord.LeaderboardChannel.ID, maxMessages, "", "", "")
			if err != nil {
				log.Errorf("failed to reattempt request" + err.Error())
			}
		} else {
			log.WithFields(log.Fields{
				string(static.UUIDKey):  ctx.Value(static.UUIDKey),
				string(static.ErrorKey): err.Error(),
			}).Error("failed to get messages from leaderboard channel")
			return
		}
	}
	for len(messages) != 0 {
		lenMsgs += len(messages)
		lastMessage := messages[len(messages)-1]
		messages, err = session.ChannelMessages(discord.LeaderboardChannel.ID, maxMessages, lastMessage.ID, "", "")
		if err != nil {
			if strings.Contains(err.Error(), "500 Internal Server Error") {
				log.Info("failed to request messages" + err.Error())
				log.Info("reattempting request...")
				time.Sleep(time.Second / 10)
				messages, err = session.ChannelMessages(discord.LeaderboardChannel.ID, maxMessages, lastMessage.ID, "", "")
				if err != nil {
					log.Errorf("failed to reattempt request" + err.Error())
				}
			} else {
				log.WithFields(log.Fields{
					string(static.UUIDKey):  ctx.Value(static.UUIDKey),
					string(static.ErrorKey): err.Error(),
				}).Error("failed to get messages from leaderboard channel")
				return
			}
		}
	}

	// send placeholder msgs if more msgs are necessary
	for i := 0; i < (nbMessagesNeeded - lenMsgs); i++ {
		_, err := session.ChannelMessageSend(discord.LeaderboardChannel.ID, "placeholder")
		if err != nil {
			log.WithFields(log.Fields{
				string(static.UUIDKey):  ctx.Value(static.UUIDKey),
				string(static.ErrorKey): err.Error(),
			}).Error("failed to send messages in leaderboard channel")
			return
		}
	}

	// delete unnecessary msgs
	for i := 0; i < (lenMsgs - nbMessagesNeeded); i++ {
		messages, err = session.ChannelMessages(discord.LeaderboardChannel.ID, 1, "", "", "")
		if err != nil {
			if strings.Contains(err.Error(), "500 Internal Server Error") {
				log.Info("failed to request messages" + err.Error())
				log.Info("reattempting request...")
				time.Sleep(time.Second / 10)
				messages, err = session.ChannelMessages(discord.LeaderboardChannel.ID, 1, "", "", "")
				if err != nil {
					log.Errorf("failed to reattempt request" + err.Error())
				}
			} else {
				log.WithFields(log.Fields{
					string(static.UUIDKey):  ctx.Value(static.UUIDKey),
					string(static.ErrorKey): err.Error(),
				}).Error("failed to get messages from leaderboard channel")
				return
			}
		}
		err := session.ChannelMessageDelete(discord.LeaderboardChannel.ID, messages[0].ID)
		if err != nil {
			log.WithFields(log.Fields{
				string(static.UUIDKey):  ctx.Value(static.UUIDKey),
				string(static.ErrorKey): err.Error(),
			}).Error("failed to delete messages in leaderboard channel")
			return
		}
	}

	// edit all current msgs
	// get all msgs
	messages, err = session.ChannelMessages(discord.LeaderboardChannel.ID, maxMessages, "", "", "")
	if err != nil {
		if strings.Contains(err.Error(), "500 Internal Server Error") {
			log.Info("failed to request messages" + err.Error())
			log.Info("reattempting request...")
			time.Sleep(time.Second / 10)
			messages, err = session.ChannelMessages(discord.LeaderboardChannel.ID, maxMessages, "", "", "")
			if err != nil {
				log.Errorf("failed to reattempt request" + err.Error())
			}
		} else {
			log.WithFields(log.Fields{
				string(static.UUIDKey):  ctx.Value(static.UUIDKey),
				string(static.ErrorKey): err.Error(),
			}).Error("failed to get messages from leaderboard channel")
			return
		}
	}
	var allMessages []*discordgo.Message
	for len(messages) != 0 {
		allMessages = append(allMessages, messages...)
		lenMsgs += len(messages)
		lastMessage := messages[len(messages)-1]
		messages, err = session.ChannelMessages(discord.LeaderboardChannel.ID, maxMessages, lastMessage.ID, "", "")
		if err != nil {
			if strings.Contains(err.Error(), "500 Internal Server Error") {
				log.Info("failed to request messages" + err.Error())
				log.Info("reattempting request...")
				time.Sleep(time.Second / 10)
				messages, err = session.ChannelMessages(discord.LeaderboardChannel.ID, maxMessages, lastMessage.ID, "", "")
				if err != nil {
					log.Errorf("failed to reattempt request" + err.Error())
				}
			} else {
				log.WithFields(log.Fields{
					string(static.UUIDKey):  ctx.Value(static.UUIDKey),
					string(static.ErrorKey): err.Error(),
				}).Error("failed to get messages from leaderboard channel")
				return
			}
		}
	}
	// reorder
	var toEditMessages []*discordgo.Message
	for n := len(allMessages) - 1; n >= 0; n-- {
		toEditMessages = append(toEditMessages, allMessages[n])
	}
	// loop edit
	contents := make([]string, nbMessagesNeeded)
	for i, player := range players {
		contents[i/5] += fmt.Sprintf("%d: %s (%d credits)\n", i+1, "<@"+player.DiscordID+">", player.Credits)
	}
	for i, msg := range toEditMessages {
		_, err := session.ChannelMessageEdit(discord.LeaderboardChannel.ID, msg.ID, contents[i])
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
