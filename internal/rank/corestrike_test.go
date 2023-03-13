package rank

import (
	"context"
	"testing"
)

func TestGetRankFromUsername(t *testing.T) {
	type args struct {
		ctx      context.Context
		username string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{name: "haashi", args: args{ctx: context.TODO(), username: "haashi"}, want: 2921},
		{name: "haashixxx", args: args{ctx: context.TODO(), username: "haashixxx"}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetRankFromUsername(tt.args.ctx, tt.args.username)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetRankFromUsername() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetRankFromUsername() = %v, want %v", got, tt.want)
			}
		})
	}
}
