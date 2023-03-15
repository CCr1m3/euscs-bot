package rank

import (
	"context"
	"testing"
)

func TestGetLeaderboardPage(t *testing.T) {
	type args struct {
		ctx    context.Context
		from   int
		size   int
		region LeaderboardRegion
	}
	tests := []struct {
		name    string
		args    args
		want    *LeaderboardResponse
		wantErr bool
	}{
		{
			name: "basictop10Europe",
			args: args{
				ctx:    context.TODO(),
				from:   1,
				size:   10,
				region: Europe,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetLeaderboardPage(tt.args.ctx, tt.args.from, tt.args.size, tt.args.region)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetLeaderboardPage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.Paging.Region != Europe {
				t.Errorf("GetLeaderboardPage().Region = %v, want %v", len(got.Players), tt.args.size)
			}
			if len(got.Players) != tt.args.size {
				t.Errorf("len(GetLeaderboardPage().Players) = %v, want %v", len(got.Players), tt.args.size)
			}
		})
	}
}
