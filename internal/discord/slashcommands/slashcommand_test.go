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
			slashcommand: Link{},
			appcommand: &discordgo.ApplicationCommand{
				Name:                     Link{}.Name(),
				Description:              Link{}.Description(),
				Options:                  Link{}.Options(),
				DefaultMemberPermissions: Link{}.RequiredPerm(),
			},
		}, want: true,
		},
		{name: "false", args: args{
			slashcommand: Link{},
			appcommand: &discordgo.ApplicationCommand{
				Name:                     Unlink{}.Name(),
				Description:              Unlink{}.Description(),
				Options:                  Unlink{}.Options(),
				DefaultMemberPermissions: Unlink{}.RequiredPerm(),
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
