package tournaments

import "github.com/euscs/euscs-bot/internal/db"

type Match struct {
	Team1  *db.Team
	Team2  *db.Team
	Winner *db.Team
}

type TournamentMode string

const (
	SingleElim TournamentMode = "SingleElim"
	DoubleElim TournamentMode = "DoubleElim"
	Swiss      TournamentMode = "Swiss"
)

type Tournament struct {
	Matches map[int]*Match
	Mode    TournamentMode
	Teams   []*db.Team
}

/*func (t Tournament) setMatchWinner(matchID int, team int) {
	var winner *db.Team
	if team == 0 {
		t.Matches[matchID].Winner = t.Matches[matchID].Team1
		winner = t.Matches[matchID].Team1
	} else {
		t.Matches[matchID].Winner = t.Matches[matchID].Team2
		winner = t.Matches[matchID].Team2
	}

}*/

func (t Tournament) getRunningMatches() []*Match {
	matches := make([]*Match, 0)
	for _, match := range t.Matches {
		if match.Team1 != nil && match.Team2 != nil && match.Winner != nil {
			matches = append(matches, match)
		}
	}
	return matches
}

func newTournament(mode TournamentMode, teams []*db.Team) Tournament {
	t := Tournament{}
	t.Mode = mode
	t.Teams = teams
	switch t.Mode {
	case SingleElim:
		pow2 := nextPowerOf2(len(teams))
		t.Matches = make(map[int]*Match, 0)
		for k := 0; k < pow2; k++ {
			t.Matches[k] = &Match{}
		}
		order := makeSeedingOrder(pow2)
		for i := 0; i < pow2; i += 2 {
			team1Index := order[i]
			team2Index := order[i+1]
			var team1 *db.Team
			var team2 *db.Team
			if team1Index < len(teams) {
				team1 = teams[team1Index]
			}
			if team2Index < len(teams) {
				team2 = teams[team2Index]
			}
			match := t.Matches[i/2]
			match.Team1 = team1
			match.Team2 = team2
			if match.Team2 == nil {
				match.Winner = team1
			}
		}

	case DoubleElim:
		pow2 := nextPowerOf2(len(teams))
		t.Matches = make(map[int]*Match, pow2*2-2)
	case Swiss:
		t.Matches = make(map[int]*Match, 0)
	}
	return t
}

func nextPowerOf2(n int) int {
	k := 1
	for k < n {
		k = k << 1
	}
	return k
}

func makeSeedingOrder(pow2 int) []int {
	ret := make([]int, 0)
	ret = append(ret, 0)
	for i := 1; i < pow2; i *= 2 {
		newRet := make([]int, 0)
		for _, val := range ret {
			newRet = append(newRet, val)
			newRet = append(newRet, i*2-1-val)
		}
		ret = newRet
	}
	return ret
}

func insert(a []int, index int, value int) []int {
	if len(a) == index { // nil or empty slice or after last element
		return append(a, value)
	}
	a = append(a[:index+1], a[index:]...) // index < len(a)
	a[index] = value
	return a
}
