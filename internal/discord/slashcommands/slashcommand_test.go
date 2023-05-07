package slashcommands

import (
	"testing"

	"github.com/bwmarrin/discordgo"
)

func Test_compareCommands(t *testing.T) {
	type args struct {
		slashcommand SlashCommand
		appcommand   *discordgo.ApplicationCommand
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "basic", args: args{
			slashcommand: Sync{},
			appcommand: &discordgo.ApplicationCommand{
				Name:                     Sync{}.Name(),
				Description:              Sync{}.Description(),
				Options:                  Sync{}.Options(),
				DefaultMemberPermissions: Sync{}.RequiredPerm(),
			},
		}, want: true,
		},
		{name: "false", args: args{
			slashcommand: Unsync{},
			appcommand: &discordgo.ApplicationCommand{
				Name:                     Unsync{}.Name(),
				Description:              Unsync{}.Description(),
				Options:                  Unsync{}.Options(),
				DefaultMemberPermissions: Unsync{}.RequiredPerm(),
			},
		}, want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := compareCommands(tt.args.slashcommand, tt.args.appcommand); got != tt.want {
				t.Errorf("compareCommands() = %v, want %v", got, tt.want)
			}
		})
	}
}
