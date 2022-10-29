package matchmaking

import (
	"errors"
	"github.com/haashi/omega-strikers-bot/internal/discord"
	"os"
	"time"

	"github.com/haashi/omega-strikers-bot/internal/db"
	"github.com/haashi/omega-strikers-bot/internal/models"
	"github.com/haashi/omega-strikers-bot/internal/rank"
	log "github.com/sirupsen/logrus"
)

func AddPlayerToQueue(playerID string, role models.Role) error {
	p, err := getOrCreatePlayer(playerID)
	if err != nil {
		return err
	}
	err = rank.UpdateRankIfNeeded(playerID)
	var tooFastErr *models.RankUpdateTooFastError
	if errors.As(err, &tooFastErr) {
	} else {
		return err
	}
	err = db.AddPlayerToQueue(p, role, int(time.Now().Unix()))
	if err != nil {
		return err
	}
	log.Infof("%s joined the queue as a %s", playerID, role)
	return nil
}

func RemovePlayerFromQueue(playerID string) error {
	p, err := getOrCreatePlayer(playerID)
	if err != nil {
		return err
	}
	err = db.RemovePlayerFromQueue(p)
	if err != nil {
		return err
	}
	log.Infof("%s left the queue", playerID)
	return nil
}

func IsPlayerInQueue(playerID string) (bool, error) {
	p, err := getOrCreatePlayer(playerID)
	if err != nil {
		return false, err
	}
	return db.IsPlayerInQueue(p)
}

func removeLongQueuers() {
	playersInQueue, _ := db.GetPlayersInQueue()
	cleanDelay := time.Hour
	if os.Getenv("mode") == "dev" {
		cleanDelay = time.Minute
	}
	for _, player := range playersInQueue {
		if time.Since(time.Unix(int64(player.EntryTime), 0)) > cleanDelay {
			db.RemovePlayerFromQueue(&player.Player)
			_, err := discord.GetSession().ChannelMessageSend(discord.AimiRequestsChannel.ID, "<@"+player.DiscordID+">, you have been removed from the queue for inactivity. Please use the /leave command next time if you didn't mean to still be in queue. If you're still here wanting to queue, /join again!")
			if err != nil {
				log.Error(err)
			}
			log.Infof("removing player %s from queue", player.OSUser)
		}
	}
}
