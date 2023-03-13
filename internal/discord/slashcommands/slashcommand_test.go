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
			slashcommand: Who{},
			appcommand: &discordgo.ApplicationCommand{
				Name:                     Who{}.Name(),
				Description:              Who{}.Description(),
				Options:                  Who{}.Options(),
				DefaultMemberPermissions: Who{}.RequiredPerm(),
			},
		}, want: true,
		},
		{name: "false", args: args{
			slashcommand: Who{},
			appcommand: &discordgo.ApplicationCommand{
				Name:                     Link{}.Name(),
				Description:              Link{}.Description(),
				Options:                  Link{}.Options(),
				DefaultMemberPermissions: Link{}.RequiredPerm(),
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
