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
	log "github.com/sirupsen/logrus"
)

func Init() {
	log.Info("starting matchmaking service")

	if os.Getenv("mode") == "dev" {
		log.Info("starting dummy players")
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
		go func() {
			for {
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
	go func() {
		for {
			deleteOldMatches()
			time.Sleep(60 * time.Second)
		}
	}()
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

func zeroFlexGoaliesSample(forwards int, flex int, goalies int, indices *[6]int) {

}

func oneFlexGoalieSample(forwards int, flex int, goalies int, indices *[6]int) {

}

func twoFlexGoaliesSample(forwards int, flex int, goalies int, indices *[6]int) {

}

func algorithm() ([]*models.Player, []*models.Player) {
	playersInQueue, _ := db.GetPlayersInQueue()
	forwards, flex, goalies := 0, 0, 0
	sort.SliceStable(playersInQueue, func(i, j int) bool { //goalie -> flex -> forward priority
		if playersInQueue[i].Role == "goalie" || playersInQueue[j].Role == "goalie" {
			return playersInQueue[i].Role == "goalie" && playersInQueue[j].Role != "goalie"
		}
		if playersInQueue[i].Role == "flex" || playersInQueue[j].Role == "flex" {
			return playersInQueue[i].Role == "flex" && playersInQueue[j].Role != "flex"
		}
		return false
	})
	for _, player := range playersInQueue {
		switch player.Role {
		case "goalie":
			goalies++
		case "flex":
			flex++
		case "forward":
			forwards++
		}
	}
	// Number of possible combinations multiplied by 4!*2/((forwards+flex-2)((forwards+flex-3))) (math don't worry about it)
	// All these formulas behave nicely and give 0 when there are no possibilities.
	zeroFlexGoalies := float64(goalies * (goalies - 1) * (forwards + flex) * (forwards + flex - 1))
	oneFlexGoalie := float64(goalies * flex * 2 * (forwards + flex - 1) * (forwards + flex - 4))
	twoFlexGoalies := float64(flex * (flex - 1) * (forwards + flex - 4) * (forwards + flex - 5))
	totalPossibilities := zeroFlexGoalies + oneFlexGoalie + twoFlexGoalies
	zeroFlexGoaliesProbability := zeroFlexGoalies / totalPossibilities
	oneOrZeroFlexGoalieProbability := oneFlexGoalie/totalPossibilities + zeroFlexGoaliesProbability
	var indices [6]int
	for i := 0; i < 1000; i++ {
		r := rand.Float64()
		if r < zeroFlexGoaliesProbability {
			zeroFlexGoaliesSample(forwards, flex, goalies, &indices)
		} else if r < oneOrZeroFlexGoalieProbability {
			oneFlexGoalieSample(forwards, flex, goalies, &indices)
		} else {
			twoFlexGoaliesSample(forwards, flex, goalies, &indices)
		}
		// TODO: check the range of the players and note the best.
	}
}
