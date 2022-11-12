package discord

import (
	"github.com/bwmarrin/discordgo"

	log "github.com/sirupsen/logrus"
)

var MatchesChannel *discordgo.Channel
var AiMiChannels *discordgo.Channel
var HowToChannel *discordgo.Channel
var AimiRequestsChannel *discordgo.Channel
var LeaderboardChannel *discordgo.Channel

func initChannels() error {
	channels, err := session.GuildChannels(GuildID)
	if err != nil {
		log.Error("failed to get guild channels: ", err.Error())
	}
	for _, channel := range channels {
		if channel.Name == "Ai.Mi" {
			AiMiChannels = channel
		}
		if channel.Name == "instructions" {
			HowToChannel = channel
		}
		if channel.Name == "matches" {
			MatchesChannel = channel
		}
		if channel.Name == "aimi-requests" {
			AimiRequestsChannel = channel
		}
		if channel.Name == "credits-leaderboard" {
			LeaderboardChannel = channel
		}
	}
	if AiMiChannels == nil {
		AiMiChannels, err = session.GuildChannelCreate(GuildID, "Ai.Mi", discordgo.ChannelTypeGuildCategory)
		if err != nil {
			log.Fatal("failed to create channel group Ai.Mi: ", err.Error())
		}
	}
	if HowToChannel == nil {
		HowToChannel, err = session.GuildChannelCreateComplex(GuildID, discordgo.GuildChannelCreateData{Name: "instructions", Type: discordgo.ChannelTypeGuildText, ParentID: AiMiChannels.ID})
		if err != nil {
			log.Fatal("failed to create channel how-to: ", err.Error())
		}
		err = session.ChannelPermissionSet(HowToChannel.ID, GuildID, discordgo.PermissionOverwriteTypeRole, 0, discordgo.PermissionSendMessages)
		if err != nil {
			log.Fatal("failed to lock channel matches: ", err.Error())
		}
		err = session.ChannelPermissionSet(HowToChannel.ID, ApplicationRole.ID, discordgo.PermissionOverwriteTypeRole, discordgo.PermissionSendMessages, 0)
		if err != nil {
			log.Fatal("failed to open channel matches for bot: ", err.Error())
		}
	}
	if MatchesChannel == nil {
		MatchesChannel, err = session.GuildChannelCreateComplex(GuildID, discordgo.GuildChannelCreateData{Name: "matches", Type: discordgo.ChannelTypeGuildText, ParentID: AiMiChannels.ID})
		if err != nil {
			log.Fatal("failed to create channel matches: ", err.Error())
		}
		err = session.ChannelPermissionSet(MatchesChannel.ID, GuildID, discordgo.PermissionOverwriteTypeRole, 0, discordgo.PermissionSendMessages)
		if err != nil {
			log.Fatal("failed to lock channel matches: ", err.Error())
		}
		err = session.ChannelPermissionSet(MatchesChannel.ID, ApplicationRole.ID, discordgo.PermissionOverwriteTypeRole, discordgo.PermissionSendMessages, 0)
		if err != nil {
			log.Fatal("failed to open channel matches for bot: ", err.Error())
		}
	}
	if AimiRequestsChannel == nil {
		AimiRequestsChannel, err = session.GuildChannelCreateComplex(GuildID, discordgo.GuildChannelCreateData{Name: "aimi-requests", Type: discordgo.ChannelTypeGuildText, ParentID: AiMiChannels.ID})
		if err != nil {
			log.Fatal("failed to create channel aimi-requests: ", err.Error())
		}
	}
	if LeaderboardChannel == nil {
		LeaderboardChannel, err = session.GuildChannelCreateComplex(GuildID, discordgo.GuildChannelCreateData{Name: "credits-leaderboard", Type: discordgo.ChannelTypeGuildText, ParentID: AiMiChannels.ID})
		if err != nil {
			log.Fatal("failed to create channel credits-leaderboard: ", err.Error())
		}
		err = session.ChannelPermissionSet(LeaderboardChannel.ID, GuildID, discordgo.PermissionOverwriteTypeRole, 0, discordgo.PermissionSendMessages)
		if err != nil {
			log.Fatal("failed to lock channel credits-leaderboard: ", err.Error())
		}
		err = session.ChannelPermissionSet(LeaderboardChannel.ID, ApplicationRole.ID, discordgo.PermissionOverwriteTypeRole, discordgo.PermissionSendMessages, 0)
		if err != nil {
			log.Fatal("failed to open channel credits-leaderboard for bot: ", err.Error())
		}
	}
	err = initHowTo()
	if err != nil {
		log.Fatal("failed to init channel how-to: ", err.Error())
	}
	return nil
}
