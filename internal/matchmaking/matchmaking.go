package matchmaking

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"sort"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	"github.com/haashi/omega-strikers-bot/internal/db"
	"github.com/haashi/omega-strikers-bot/internal/discord"
	"github.com/haashi/omega-strikers-bot/internal/env"
	"github.com/haashi/omega-strikers-bot/internal/models"
	"github.com/haashi/omega-strikers-bot/internal/scheduled"
	"github.com/haashi/omega-strikers-bot/internal/utils"
	log "github.com/sirupsen/logrus"
)

func Init() {
	log.Info("starting matchmaking service")
	if env.Mode == env.DEV {
		ctx := context.TODO()
		log.Debug("starting dummy players")
		dummies := make([]string, 0)
		dummiesUsername := [30]string{"BaluGoalie", "Haaashi", "Piols", "Balu", "Lynx_", "Connax", "Masus", "Kolashiu", "Buntaoo", "Ballgrabber", "Jimray3", "IamTrusty", "Kidpan", "MathieuCalip", "Goku", "czem", "HHaie KuKi", "KeeperofLolis", "Madoushy", "LDC", "Yatta", "Immaculator", "goalkeeper diff", "Funii", "kirby", "mascha", "Thezs", "Cognity", "Sm1le", "Yuume"}
		r := rand.New(rand.NewSource(2))
		for i := 0; i < 30; i++ {
			playerID := fmt.Sprintf("%d", r.Intn(math.MaxInt64))
			player, err := db.GetOrCreatePlayerById(ctx, playerID)
			if err != nil {
				log.Error(err)
			}
			player.OSUser = strings.ToLower(dummiesUsername[i])
			err = db.UpdatePlayer(ctx, player)
			if err != nil {
				log.Error(err)
			}
			dummies = append(dummies, playerID)
		}
		dummiesFunc := func() {
			ctx := context.WithValue(context.Background(), models.UUIDKey, uuid.New())
			playerID := dummies[rand.Intn(len(dummies))]
			player, _ := db.GetOrCreatePlayerById(ctx, playerID)
			roles := make([]models.Role, 0)
			roles = append(roles,
				models.RoleGoalie,
				models.RoleFlex,
				models.RoleForward,
				models.RoleForward,
				models.RoleForward)
			inMatch, _ := IsPlayerInMatch(ctx, player.DiscordID)
			inQueue, _ := IsPlayerInQueue(ctx, player.DiscordID)
			if !inQueue && !inMatch {
				err := AddPlayerToQueue(ctx, player.DiscordID, roles[rand.Intn(len(roles))])
				if err != nil {
					log.Error(err)
				}
			}
		}
		scheduled.TaskManager.Add(scheduled.Task{ID: "dummies", Run: dummiesFunc, Frequency: time.Second})
	}
	scheduled.TaskManager.Add(scheduled.Task{ID: "updatesession", Run: updateStatus, Frequency: time.Second * 5})
	scheduled.TaskManager.Add(scheduled.Task{ID: "trycreatingmatch", Run: tryCreatingMatch, Frequency: time.Second * 15})
	scheduled.TaskManager.Add(scheduled.Task{ID: "closeoldmatches", Run: deleteOldMatches, Frequency: time.Minute})
	scheduled.TaskManager.Add(scheduled.Task{ID: "removelongqueuers", Run: removeLongQueuers, Frequency: time.Minute})
	scheduled.TaskManager.Add(scheduled.Task{ID: "threadcleanup", Run: threadCleanUp, Frequency: time.Hour})
	ctx := context.WithValue(context.Background(), models.UUIDKey, uuid.New())
	waitingForVoteMatches, err := db.GetWaitingForVotesMatches(ctx)
	if err != nil {
		log.Error("failed to get matches with a vote in progress:" + err.Error())
	} else {
		if len(waitingForVoteMatches) > 0 {
			for _, match := range waitingForVoteMatches {
				scheduled.TaskManager.Add(scheduled.Task{ID: "matchvote" + match.ID, Frequency: time.Second, Run: func() { handleMatchVoteResult(match) }})
			}
		}
	}
}

func updateStatus() {
	ctx := context.WithValue(context.Background(), models.UUIDKey, uuid.New())
	session := discord.GetSession()
	playersInQueue, _ := db.GetPlayersInQueue(ctx)
	goaliesCount, err := db.GetGoaliesCountInQueue(ctx)
	if err != nil {
		log.WithFields(log.Fields{
			string(models.UUIDKey):  ctx.Value(models.UUIDKey),
			string(models.ErrorKey): err.Error(),
		}).Error("failed to get goalies count")
		return
	}
	forwardsCount, err := db.GetForwardsCountInQueue(ctx)
	if err != nil {
		log.WithFields(log.Fields{
			string(models.UUIDKey):  ctx.Value(models.UUIDKey),
			string(models.ErrorKey): err.Error(),
		}).Error("failed to get forwards count")
		return
	}
	queueSize := len(playersInQueue)
	var act []*discordgo.Activity
	act = append(act, &discordgo.Activity{Name: fmt.Sprintf("%d people queuing (%d goalies, %d forwards)", queueSize, goaliesCount, forwardsCount), Type: discordgo.ActivityTypeWatching})
	err = session.UpdateStatusComplex(discordgo.UpdateStatusData{Activities: act})
	if err != nil {
		log.Error(err)
	}
}

func tryCreatingMatch() {
	ctx := context.WithValue(context.Background(), models.UUIDKey, uuid.New())
	playersInQueue, err := db.GetPlayersInQueue(ctx)
	if err != nil {
		log.Error(err)
	}
	goalieInQueue, err := db.GetGoaliesCountInQueue(ctx)
	if err != nil {
		log.Error(err)
	}
	forwardInQueue, err := db.GetForwardsCountInQueue(ctx)
	if err != nil {
		log.Error(err)
	}
	if len(playersInQueue) >= 6 && goalieInQueue >= 2 && forwardInQueue >= 4 {
		team1, team2 := algorithm(ctx)
		if len(team1) == 0 {
			log.Debug("match not created, algorithm deemed no match of quality")
			return
		}
		err := createNewMatch(ctx, team1, team2)
		if err != nil {
			log.Error("could not create new match: ", err)
		} else {
			players := append(team1, team2...)
			for _, player := range players {
				err := RemovePlayerFromQueue(ctx, player.DiscordID)
				if err != nil {
					log.Error("could not make player leave queue: ", err)
				}
			}
		}
	} else {
		return
	}
}

// Haashi please don't kill me, I'm just optimizing. This needs to be as fast as possible.
func zeroFlexGoaliesSample(forwards int, flex int, goalies int) [6]int {
	indices := [6]int{}
	indices[0] = utils.FastRandN(goalies)
	indices[1] = utils.FastRandN(goalies - 1)
	if indices[1] >= indices[0] {
		indices[1]++
	}
	indices[2] = utils.FastRandN(forwards+flex) + goalies
	indices[3] = utils.FastRandN(forwards+flex-1) + goalies
	if indices[3] >= indices[2] {
		indices[3]++
	}
	indices[4] = utils.FastRandN(forwards+flex) + goalies
	for indices[4] == indices[3] || indices[4] == indices[2] {
		indices[4] = utils.FastRandN(forwards+flex) + goalies
	}
	indices[5] = utils.FastRandN(forwards+flex) + goalies
	for indices[5] == indices[4] || indices[5] == indices[3] || indices[5] == indices[2] {
		indices[5] = utils.FastRandN(forwards+flex) + goalies
	}
	return indices
}

func oneFlexGoalieSample(forwards int, flex int, goalies int) [6]int {
	indices := [6]int{}
	indices[0] = utils.FastRandN(goalies)
	indices[1] = utils.FastRandN(flex) + goalies
	indices[2] = utils.FastRandN(flex+forwards-1) + goalies
	if indices[2] >= indices[1] {
		indices[2]++
	}
	indices[3] = utils.FastRandN(forwards+flex) + goalies
	for indices[3] == indices[2] || indices[3] == indices[1] {
		indices[3] = utils.FastRandN(forwards+flex) + goalies
	}
	indices[4] = utils.FastRandN(forwards+flex) + goalies
	for indices[4] == indices[3] || indices[4] == indices[2] || indices[4] == indices[1] {
		indices[4] = utils.FastRandN(forwards+flex) + goalies
	}
	indices[5] = utils.FastRandN(forwards+flex) + goalies
	for indices[5] == indices[4] || indices[5] == indices[3] || indices[5] == indices[2] || indices[5] == indices[1] {
		indices[5] = utils.FastRandN(forwards+flex) + goalies
	}
	return indices
}

func twoFlexGoaliesSample(forwards int, flex int, goalies int) [6]int {
	indices := [6]int{}
	indices[0] = utils.FastRandN(flex) + goalies
	indices[1] = utils.FastRandN(flex-1) + goalies
	if indices[1] >= indices[0] {
		indices[1]++
	}
	indices[2] = utils.FastRandN(forwards+flex) + goalies
	for indices[2] == indices[1] || indices[2] == indices[0] {
		indices[2] = utils.FastRandN(forwards+flex) + goalies
	}
	indices[3] = utils.FastRandN(forwards+flex) + goalies
	for indices[3] == indices[2] || indices[3] == indices[1] || indices[3] == indices[0] {
		indices[3] = utils.FastRandN(forwards+flex) + goalies
	}
	indices[4] = utils.FastRandN(forwards+flex) + goalies
	for indices[4] == indices[3] || indices[4] == indices[2] || indices[4] == indices[1] || indices[4] == indices[0] {
		indices[4] = utils.FastRandN(forwards+flex) + goalies
	}
	indices[5] = utils.FastRandN(forwards+flex) + goalies
	for indices[5] == indices[4] || indices[5] == indices[3] || indices[5] == indices[2] || indices[5] == indices[1] || indices[5] == indices[0] {
		indices[5] = utils.FastRandN(forwards+flex) + goalies
	}
	return indices
}

func evaluatePlayers(indices *[6]int, players []*models.QueuedPlayer) float64 {
	const eloRange = 500
	maxElo, minElo := -1, 1<<20
	sumQueueTime := 0.0
	for i := 0; i < 6; i++ {
		player := players[indices[i]]
		if player.Elo > maxElo {
			maxElo = player.Elo
		}
		if player.Elo < minElo {
			minElo = player.Elo
		}
		sumQueueTime += time.Since(time.Unix(int64(player.EntryTime), 0)).Minutes()
	}
	softMaxElo := maxElo
	// This is "temporary", as in for the foreseeable future. Makes all omega players eligible to play with all other omegas.
	if softMaxElo > 3400 {
		softMaxElo = 3400
	}
	if softMaxElo-minElo > 500 { // This is negative therefore match won't be accepted
		return float64(eloRange - (maxElo - minElo))
	}
	return float64(eloRange-(softMaxElo-minElo)) + sumQueueTime*10.0
}

func evaluateTeams(team1 []*models.Player, team2 []*models.Player) float64 {
	return float64(team1[0].Elo)*0.4 + float64(team1[1].Elo)*0.3 + float64(team1[2].Elo)*0.3 - (float64(team2[0].Elo)*0.4 + float64(team2[1].Elo)*0.3 + float64(team2[2].Elo)*0.3)
}

func balanceTeams(indices *[6]int, players []*models.QueuedPlayer) ([]*models.Player, []*models.Player) {
	fwdsSplit := [6][4]int{{1, 2, 3, 4}, {1, 3, 2, 4}, {1, 4, 2, 3}, {2, 3, 1, 4}, {2, 4, 1, 3}, {3, 4, 1, 2}}
	bestSplit := fwdsSplit[0]
	bestBalance := float64(1 << 20)
	for _, split := range fwdsSplit {
		team1 := []*models.Player{&players[indices[0]].Player, &players[indices[split[0]+1]].Player, &players[indices[split[1]+1]].Player}
		team2 := []*models.Player{&players[indices[1]].Player, &players[indices[split[2]+1]].Player, &players[indices[split[3]+1]].Player}
		balance := evaluateTeams(team1, team2)
		if math.Abs(balance) < math.Abs(bestBalance) {
			bestSplit = split
			bestBalance = balance
		}
	}
	team1 := []*models.Player{&players[indices[0]].Player, &players[indices[bestSplit[0]+1]].Player, &players[indices[bestSplit[1]+1]].Player}
	team2 := []*models.Player{&players[indices[1]].Player, &players[indices[bestSplit[2]+1]].Player, &players[indices[bestSplit[3]+1]].Player}
	log.Debugf("elos: 1st team: %d %d %d, 2nd team: %d %d %d", players[indices[0]].Elo, players[indices[bestSplit[0]+1]].Elo, players[indices[bestSplit[1]+1]].Elo, players[indices[1]].Elo, players[indices[bestSplit[2]+1]].Elo, players[indices[bestSplit[3]+1]].Elo)
	log.Debugf("best team balance: %.2f", bestBalance)
	return team1, team2
}

func algorithm(ctx context.Context) ([]*models.Player, []*models.Player) {
	playersInQueue, err := db.GetPlayersInQueue(ctx)
	if err != nil {
		log.Error(err)
	}
	forwards, flex, goalies := 0, 0, 0
	sort.SliceStable(playersInQueue, func(i, j int) bool { //goalie -> flex -> forward priority
		return (playersInQueue[i].Role == "goalie" && playersInQueue[j].Role != "goalie") || (playersInQueue[i].Role == "flex" && playersInQueue[j].Role == "forward")
	})
	log.Debug("players in sorted queue:")
	for i, player := range playersInQueue {
		log.Debugf("%d %s %s %d", i, player.Role, player.OSUser, player.Elo)
	}
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
	log.Debugf("relative probabilities for X flex goalies: 0: %.2f, 1: %.2f, 2: %.2f", zeroFlexGoalies, oneFlexGoalie, twoFlexGoalies)
	var indices [6]int
	var bestIndices [6]int
	bestQuality := -1.0
	samplesTaken := 1000
	if env.Mode == env.DEV {
		samplesTaken = 100
	}
	for i := 0; i < samplesTaken; i++ {
		r := rand.Float64()
		if r < zeroFlexGoaliesProbability {
			indices = zeroFlexGoaliesSample(forwards, flex, goalies)
		} else if r < oneOrZeroFlexGoalieProbability {
			indices = oneFlexGoalieSample(forwards, flex, goalies)
		} else {
			indices = twoFlexGoaliesSample(forwards, flex, goalies)
		}
		quality := evaluatePlayers(&indices, playersInQueue)
		if quality > bestQuality {
			bestQuality, bestIndices = quality, indices
		}
	}
	if bestQuality < 0 {
		return []*models.Player{}, []*models.Player{}
	}
	log.Debugf("best indices: %d %d %d %d %d %d", bestIndices[0], bestIndices[1], bestIndices[2], bestIndices[3], bestIndices[4], bestIndices[5])
	return balanceTeams(&bestIndices, playersInQueue)
}
