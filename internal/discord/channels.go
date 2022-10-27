package discord

import (
	"github.com/bwmarrin/discordgo"

	log "github.com/sirupsen/logrus"
)

var MatchesChannel *discordgo.Channel
var AiMiChannels *discordgo.Channel
var HowToChannel *discordgo.Channel

func initChannels() error {
	channels, err := session.GuildChannels(GuildID)
	if err != nil {
		log.Errorf("failed to get guild channels")
	}
	for _, channel := range channels {
		if channel.Name == "Ai.Mi" {
			AiMiChannels = channel
		}
		if channel.Name == "how-to" {
			HowToChannel = channel
		}
		if channel.Name == "matches" {
			MatchesChannel = channel
		}
	}
	if AiMiChannels == nil {
		AiMiChannels, err = session.GuildChannelCreate(GuildID, "Ai.Mi", discordgo.ChannelTypeGuildCategory)
		if err != nil {
			log.Fatalf("failed to create channel group Ai.Mi")
		}
	}
	if HowToChannel == nil {
		HowToChannel, err = session.GuildChannelCreateComplex(GuildID, discordgo.GuildChannelCreateData{Name: "how-to", Type: discordgo.ChannelTypeGuildText, ParentID: AiMiChannels.ID})
		if err != nil {
			log.Fatalf("failed to create channel how-to")
		}
		err = session.ChannelPermissionSet(HowToChannel.ID, GuildID, discordgo.PermissionOverwriteTypeRole, 0, discordgo.PermissionSendMessages)
		if err != nil {
			log.Fatalf("failed to lock channel matches")
		}
	}
	if MatchesChannel == nil {
		MatchesChannel, err = session.GuildChannelCreateComplex(GuildID, discordgo.GuildChannelCreateData{Name: "matches", Type: discordgo.ChannelTypeGuildText, ParentID: AiMiChannels.ID})
		if err != nil {
			log.Fatalf("failed to create channel matches")
		}
		err = session.ChannelPermissionSet(MatchesChannel.ID, GuildID, discordgo.PermissionOverwriteTypeRole, 0, discordgo.PermissionSendMessages)
		if err != nil {
			log.Fatalf("failed to lock channel matches")
		}
		err = session.ChannelPermissionSet(MatchesChannel.ID, ApplicationRole.ID, discordgo.PermissionOverwriteTypeRole, discordgo.PermissionSendMessages, 0)
		if err != nil {
			log.Fatalf("failed to open channel matches for bot")
		}
	}
	return nil
}
