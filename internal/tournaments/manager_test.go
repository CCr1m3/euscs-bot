package tournaments

import (
	"reflect"
	"testing"

	"github.com/euscs/euscs-bot/internal/db"
	"github.com/google/go-cmp/cmp"
)

func Test_makeSeedingOrder(t *testing.T) {
	type args struct {
		pow2 int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{name: "2", args: args{pow2: 2}, want: []int{0, 1}},
		{name: "4", args: args{pow2: 4}, want: []int{0, 3, 1, 2}},
		{name: "8", args: args{pow2: 8}, want: []int{0, 7, 3, 4, 1, 6, 2, 5}},
		{name: "16", args: args{pow2: 16}, want: []int{0, 15, 7, 8, 3, 12, 4, 11, 1, 14, 6, 9, 2, 13, 5, 10}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := makeSeedingOrder(tt.args.pow2); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("makeSeedingOrder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newTournament(t *testing.T) {
	team1 := &db.Team{Name: "team1", Players: db.Players{{}, {}, {}}}
	team2 := &db.Team{Name: "team2", Players: db.Players{{}, {}, {}}}
	team3 := &db.Team{Name: "team3", Players: db.Players{{}, {}, {}}}
	team4 := &db.Team{Name: "team4", Players: db.Players{{}, {}, {}}}
	team5 := &db.Team{Name: "team5", Players: db.Players{{}, {}, {}}}
	team6 := &db.Team{Name: "team6", Players: db.Players{{}, {}, {}}}
	team7 := &db.Team{Name: "team7", Players: db.Players{{}, {}, {}}}
	team8 := &db.Team{Name: "team8", Players: db.Players{{}, {}, {}}}
	type args struct {
		mode  TournamentMode
		teams []*db.Team
	}
	tests := []struct {
		name string
		args args
		want Tournament
	}{
		{
			name: "singleelim", args: args{
				mode: SingleElim,
				teams: []*db.Team{
					team1, team2, team3, team4, team5, team6, team7, team8,
				},
			}, want: Tournament{
				Mode: SingleElim,
				Teams: []*db.Team{
					team1, team2, team3, team4, team5, team6, team7, team8,
				},
				Matches: map[int]*Match{
					0: {Team1: team1, Team2: team8, Winner: nil},
					1: {Team1: team4, Team2: team5, Winner: nil},
					2: {Team1: team2, Team2: team7, Winner: nil},
					3: {Team1: team3, Team2: team6, Winner: nil},
					4: {Team1: nil, Team2: nil, Winner: nil},
					5: {Team1: nil, Team2: nil, Winner: nil},
					6: {Team1: nil, Team2: nil, Winner: nil},
					7: {Team1: nil, Team2: nil, Winner: nil},
				},
			},
		},
		{name: "singleelimwithBYES", args: args{
			mode: SingleElim,
			teams: []*db.Team{
				team1, team2, team3, team4, team5, team6,
			},
		}, want: Tournament{
			Mode: SingleElim,
			Teams: []*db.Team{
				team1, team2, team3, team4, team5, team6,
			},
			Matches: map[int]*Match{
				0: {Team1: team1, Team2: nil, Winner: team1},
				1: {Team1: team4, Team2: team5, Winner: nil},
				2: {Team1: team2, Team2: nil, Winner: team2},
				3: {Team1: team3, Team2: team6, Winner: nil},
				4: {Team1: nil, Team2: nil, Winner: nil},
				5: {Team1: nil, Team2: nil, Winner: nil},
				6: {Team1: nil, Team2: nil, Winner: nil},
				7: {Team1: nil, Team2: nil, Winner: nil},
			},
		},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newTournament(tt.args.mode, tt.args.teams); !cmp.Equal(tt.want, got) {
				t.Errorf("tournaments are different: %s", cmp.Diff(tt.want, got))
			}
		})
	}
}
