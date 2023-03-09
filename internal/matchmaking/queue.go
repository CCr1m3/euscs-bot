package matchmaking

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/haashi/omega-strikers-bot/internal/discord"
	"github.com/haashi/omega-strikers-bot/internal/env"

	"github.com/haashi/omega-strikers-bot/internal/db"
	"github.com/haashi/omega-strikers-bot/internal/models"
	"github.com/haashi/omega-strikers-bot/internal/rank"
	log "github.com/sirupsen/logrus"
)

func AddPlayerToQueue(ctx context.Context, playerID string, role models.Role) error {
	p, err := db.GetOrCreatePlayerById(ctx, playerID)
	if err != nil {
		log.WithFields(log.Fields{
			string(models.UUIDKey):     ctx.Value(models.UUIDKey),
			string(models.PlayerIDKey): playerID,
			string(models.ErrorKey):    err.Error(),
		}).Error("failed to get or create player")
		return err
	}
	err = rank.UpdateRankIfNeeded(ctx, playerID)
	if err != nil {
		var tooFastErr *models.RankUpdateTooFastError
		if errors.As(err, &tooFastErr) {
		} else {
			log.WithFields(log.Fields{
				string(models.UUIDKey):     ctx.Value(models.UUIDKey),
				string(models.PlayerIDKey): playerID,
				string(models.ErrorKey):    err.Error(),
			}).Error("failed to update player rank before joining queue")
		}
	}

	err = db.AddPlayerToQueue(ctx, p, role, int(time.Now().Unix()))
	if err != nil {
		log.WithFields(log.Fields{
			string(models.UUIDKey):     ctx.Value(models.UUIDKey),
			string(models.PlayerIDKey): playerID,
			string(models.ErrorKey):    err.Error(),
		}).Error("failed to add player to queue")
		return err
	}
	log.WithFields(log.Fields{
		string(models.UUIDKey):      ctx.Value(models.UUIDKey),
		string(models.PlayerIDKey):  playerID,
		string(models.QueueRoleKey): role,
	}).Info("player joined the queue")
	return nil
}

func RemovePlayerFromQueue(ctx context.Context, playerID string) error {
	p, err := db.GetOrCreatePlayerById(ctx, playerID)
	if err != nil {
		log.WithFields(log.Fields{
			string(models.UUIDKey):     ctx.Value(models.UUIDKey),
			string(models.PlayerIDKey): playerID,
			string(models.ErrorKey):    err.Error(),
		}).Error("failed to get or create player")
		return err
	}
	err = db.RemovePlayerFromQueue(ctx, p)
	if err != nil {
		return err
	}
	log.Infof("%s left the queue", playerID)
	return nil
}

func IsPlayerInQueue(ctx context.Context, playerID string) (bool, error) {
	p, err := db.GetOrCreatePlayerById(ctx, playerID)
	if err != nil {
		log.WithFields(log.Fields{
			string(models.UUIDKey):     ctx.Value(models.UUIDKey),
			string(models.PlayerIDKey): playerID,
			string(models.ErrorKey):    err.Error(),
		}).Error("failed to get or create player")
		return false, err
	}
	res, err := db.IsPlayerInQueue(ctx, p)
	if err != nil {
		log.WithFields(log.Fields{
			string(models.UUIDKey):     ctx.Value(models.UUIDKey),
			string(models.PlayerIDKey): playerID,
			string(models.ErrorKey):    err.Error(),
		}).Error("failed to check if player is in queue")
		return false, err
	}
	return res, nil
}

func removeLongQueuers() {
	ctx := context.WithValue(context.Background(), models.UUIDKey, uuid.New())
	playersInQueue, err := db.GetPlayersInQueue(ctx)
	if err != nil {
		log.Error(err)
		return
	}
	cleanDelay := time.Hour
	if env.Mode == env.DEV {
		cleanDelay = time.Minute
	}
	for _, player := range playersInQueue {
		if time.Since(time.Unix(int64(player.EntryTime), 0)) > cleanDelay {
			log.Infof("removing player %s from queue", player.OSUser)
			err = db.RemovePlayerFromQueue(ctx, &player.Player)
			if err != nil {
				log.Error(err)
				continue
			}
			_, err := discord.GetSession().ChannelMessageSend(discord.AimiRequestsChannel.ID, "<@"+player.DiscordID+">, you have been removed from the queue for inactivity. Please use the /leave command next time if you didn't mean to still be in queue. If you're still here wanting to queue, /join again!")
			if err != nil {
				log.Error(err)
			}
		}
	}
}
