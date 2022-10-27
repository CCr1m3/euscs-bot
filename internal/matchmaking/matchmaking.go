package matchmaking

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/haashi/omega-strikers-bot/internal/db"
	"github.com/haashi/omega-strikers-bot/internal/discord"
	"github.com/haashi/omega-strikers-bot/internal/models"
	"github.com/haashi/omega-strikers-bot/internal/scheduled"
	log "github.com/sirupsen/logrus"
)

func Init() {
	log.Info("starting matchmaking service")
	if os.Getenv("mode") == "dev" {
		log.Debug("starting dummy players")
		dummies := make([]string, 0)
		dummiesUsername := [30]string{"BaluGoalie", "Haaashi", "Piols", "Balu", "Lynx_", "Connax", "Masus", "Kolashiu", "Buntaoo", "Ballgrabber", "Jimray3", "IamTrusty", "Kidpan", "MathieuCalip", "Goku", "czem", "HHaie KuKi", "KeeperofLolis", "Madoushy", "LDC", "Yatta", "Immaculator", "goalkeeper diff", "Funii", "kirby", "mascha", "Thezs", "Cognity", "Sm1le", "Yuume"}
		r := rand.New(rand.NewSource(2))
		for i := 0; i < 30; i++ {
			playerID := fmt.Sprintf("%d", r.Intn(math.MaxInt64))
			player, err := getOrCreatePlayer(playerID)
			if err != nil {
				log.Error(err)
			}
			player.OSUser = strings.ToLower(dummiesUsername[i])
			err = db.UpdatePlayer(player)
			if err != nil {
				log.Error(err)
			}
			dummies = append(dummies, playerID)
		}
		dummiesFunc := func() {
			playerID := dummies[rand.Intn(len(dummies))]
			player, _ := getOrCreatePlayer(playerID)
			roles := make([]models.Role, 0)
			roles = append(roles,
				models.RoleGoalie,
				models.RoleFlex,
				models.RoleForward,
				models.RoleForward,
				models.RoleForward)
			inMatch, _ := IsPlayerInMatch(player.DiscordID)
			inQueue, _ := IsPlayerInQueue(player.DiscordID)
			if !inQueue && !inMatch {
				err := AddPlayerToQueue(player.DiscordID, roles[rand.Intn(len(roles))])
				if err != nil {
					log.Error(err)
				}
			}
		}
		scheduled.TaskManager.Add(scheduled.Task{ID: "dummies", Run: dummiesFunc, Frequency: time.Second})
	}
	scheduled.TaskManager.Add(scheduled.Task{ID: "updatesession", Run: updateStatus, Frequency: time.Second * 15})
	scheduled.TaskManager.Add(scheduled.Task{ID: "trycreatingmatch", Run: tryCreatingMatch, Frequency: time.Second * 15})
	scheduled.TaskManager.Add(scheduled.Task{ID: "closeoldmatches", Run: deleteOldMatches, Frequency: time.Minute})
}

func updateStatus() {
	session := discord.GetSession()
	playersInQueue, _ := db.GetPlayersInQueue()
	queueSize := len(playersInQueue)
	var act []*discordgo.Activity
	act = append(act, &discordgo.Activity{Name: fmt.Sprintf("%d people queuing", queueSize), Type: discordgo.ActivityTypeWatching})
	err := session.UpdateStatusComplex(discordgo.UpdateStatusData{Activities: act})
	if err != nil {
		log.Error(err)
	}
}

func tryCreatingMatch() {
	playersInQueue, _ := db.GetPlayersInQueue()
	goalieInQueue, err := db.GetGoaliesCountInQueue()
	if err != nil {
		log.Error(err)
	}
	forwardInQueue, err := db.GetForwardsCountInQueue()
	if err != nil {
		log.Error(err)
	}
	if len(playersInQueue) >= 6 && goalieInQueue >= 2 && forwardInQueue >= 4 {
		team1, team2 := algorithm()

		err := createNewMatch(team1, team2)
		if err != nil {
			log.Error("could not create new match: ", err)
		} else {
			players := append(team1, team2...)
			for _, player := range players {
				err := RemovePlayerFromQueue(player.DiscordID)
				if err != nil {
					log.Error("could not make player leave queue: ", err)
				}
			}
		}
	} else {
		return
	}
}

func algorithm() ([]*models.Player, []*models.Player) {
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
	forwards := make([]*models.QueuedPlayer, 0)
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
	team1 := []*models.Player{&goalie1.Player}
	team2 := []*models.Player{&goalie2.Player}
	team1 = append(team1, &forwards[0].Player, &forwards[1].Player)
	team2 = append(team2, &forwards[2].Player, &forwards[3].Player)
	return team1, team2
}
