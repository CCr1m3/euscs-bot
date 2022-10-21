package matchmaking

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"sort"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/haashi/omega-strikers-bot/internal/db"
	"github.com/haashi/omega-strikers-bot/internal/discord"
	"github.com/haashi/omega-strikers-bot/internal/match"
	log "github.com/sirupsen/logrus"
)

var playersPerGame int = 6

func Init() {
	log.Info("starting matchmaking service")

	if os.Getenv("mode") == "dev" {
		dummies := make([]string, 0)
		for i := 0; i < 30; i++ {
			playerID := fmt.Sprintf("%d", rand.Intn(math.MaxInt64))
			db.CreatePlayer(playerID)
			dummies = append(dummies, playerID)
		}
		go func() {
			for {
				playerID := dummies[rand.Intn(len(dummies))]
				player, _ := db.GetPlayer(playerID)
				roles := make([]string, 0)
				roles = append(roles,
					"goalie",
					"flex",
					"forward",
					"forward",
					"forward")
				inMatch, _ := player.IsInMatch()
				if !player.IsInQueue() && !inMatch {
					player.AddToQueue(roles[rand.Intn(len(roles))])
					time.Sleep(2 * time.Second)
				}
			}
		}()
	}
	go func() {
		for {
			session := discord.GetSession()
			playersInQueue, _ := db.GetPlayersInQueue()
			queueSize := len(playersInQueue)
			var act []*discordgo.Activity
			act = append(act, &discordgo.Activity{Name: fmt.Sprintf("%d people queuing", queueSize), Type: discordgo.ActivityTypeWatching})
			err := session.UpdateStatusComplex(discordgo.UpdateStatusData{Activities: act})
			if err != nil {
				log.Error(err)
			}
			time.Sleep(15 * time.Second)

			tryCreatingMatch()
		}
	}()
}

func tryCreatingMatch() {
	log.Info("trying to create a match")
	playersInQueue, _ := db.GetPlayersInQueue()
	goalieInQueue, err := db.GetGoaliesCountInQueue()
	if err != nil {
		log.Error(err)
	}
	forwardInQueue, err := db.GetForwardsCountInQueue()
	if err != nil {
		log.Error(err)
	}
	if len(playersInQueue) >= playersPerGame && goalieInQueue >= 2 && forwardInQueue >= 4 {
		team1, team2 := algorithm()
		players := append(team1, team2...)
		for _, player := range players {
			player.LeaveQueue()
		}
		log.Info(fmt.Sprintf("enough people (%d) or goalie (%d) or forward (%d) in queue", len(playersInQueue), goalieInQueue, forwardInQueue))
		err := match.New(team1, team2)
		if err != nil {
			log.Errorf("match creation went wrong")
		}
	} else {
		log.Info(fmt.Sprintf("not enough people (%d) or goalie (%d) or forward (%d) in queue", len(playersInQueue), goalieInQueue, forwardInQueue))
		return
	}
}

func algorithm() ([]*db.Player, []*db.Player) {
	playersInQueue, _ := db.GetPlayersInQueue()
	rand.Shuffle(len(playersInQueue), func(i, j int) { playersInQueue[i], playersInQueue[j] = playersInQueue[j], playersInQueue[i] })
	sort.SliceStable(playersInQueue, func(i, j int) bool { //goalie priority
		if playersInQueue[i].Role == "goalie" && playersInQueue[j].Role != "goalie" {
			return true
		}
		return false
	})
	goalie1 := playersInQueue[0]
	goalie2 := playersInQueue[1]
	forwards := make([]*db.Player, 0)
	for _, player := range playersInQueue {
		if player.DiscordID == goalie1.DiscordID || player.DiscordID == goalie2.DiscordID {
			continue
		}
		if player.Role == "goalie" {
			continue
		}
		forwards = append(forwards, player)
		if len(forwards) >= 4 {
			break
		}
	}
	team1 := []*db.Player{goalie1}
	team2 := []*db.Player{goalie2}
	team1 = append(team1, forwards[0:2]...)
	team2 = append(team2, forwards[2:4]...)
	return team1, team2
}
