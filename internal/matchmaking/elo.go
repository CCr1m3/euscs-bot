package matchmaking

import "math"

const K float64 = 30

func probability(elo1 int, elo2 int) float64 {
	return 1 / (1 + math.Pow(10, float64(elo1-elo2)/400))
}

func eloChanges(team1elo int, team2elo int, team1won bool) (int, int) {
	p1 := probability(team2elo, team1elo)
	p2 := 1 - p1
	if team1won {
		return int(math.Ceil(K * (1 - p1))), int(math.Ceil(K * (0 - p2)))
	} else {
		return int(math.Ceil(K * (0 - p1))), int(math.Ceil(K * (1 - p2)))
	}
}
