package discord

import (
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

var RoleOmega *discordgo.Role
var RoleChallenger *discordgo.Role
var RoleDiamond *discordgo.Role
var RolePlatinum *discordgo.Role
var RoleGold *discordgo.Role
var RoleSilver *discordgo.Role
var RoleBronze *discordgo.Role
var RoleRookie *discordgo.Role

var RankRoles []*discordgo.Role

func initRoles() error {
	roles, err := session.GuildRoles(GuildID)
	if err != nil {
		return err
	}
	for _, role := range roles {
		if role.Name == "Omega" {
			RoleOmega = role
		}
		if role.Name == "Challenger" {
			RoleChallenger = role
		}
		if role.Name == "Diamond" {
			RoleDiamond = role
		}
		if role.Name == "Platinum" {
			RolePlatinum = role
		}
		if role.Name == "Gold" {
			RoleGold = role
		}
		if role.Name == "Silver" {
			RoleSilver = role
		}
		if role.Name == "Bronze" {
			RoleBronze = role
		}
		if role.Name == "Rookie" {
			RoleRookie = role
		}
	}
	mentionnable := true
	hoist := true
	if RoleOmega == nil {
		color := 15548997
		RoleOmega, err = session.GuildRoleCreate(GuildID, &discordgo.RoleParams{Name: "Omega", Color: &color, Mentionable: &mentionnable, Hoist: &hoist})
		if err != nil {
			log.Fatalf("failed to create role RoleOmega")
		}
	}
	if RoleChallenger == nil {
		color := 10181046
		RoleChallenger, err = session.GuildRoleCreate(GuildID, &discordgo.RoleParams{Name: "Challenger", Color: &color, Mentionable: &mentionnable, Hoist: &hoist})
		if err != nil {
			log.Fatalf("failed to create role RoleChallenger")
		}
	}
	if RoleDiamond == nil {
		color := 3447003
		RoleDiamond, err = session.GuildRoleCreate(GuildID, &discordgo.RoleParams{Name: "Diamond", Color: &color, Mentionable: &mentionnable, Hoist: &hoist})
		if err != nil {
			log.Fatalf("failed to create role RoleDiamond")
		}
	}
	if RolePlatinum == nil {
		color := 2067276
		RolePlatinum, err = session.GuildRoleCreate(GuildID, &discordgo.RoleParams{Name: "Platinum", Color: &color, Mentionable: &mentionnable, Hoist: &hoist})
		if err != nil {
			log.Fatalf("failed to create role RolePlatinum")
		}
	}
	if RoleGold == nil {
		color := 16776960
		RoleGold, err = session.GuildRoleCreate(GuildID, &discordgo.RoleParams{Name: "Gold", Color: &color, Mentionable: &mentionnable, Hoist: &hoist})
		if err != nil {
			log.Fatalf("failed to create role RoleGold")
		}
	}
	if RoleSilver == nil {
		color := 12370112
		RoleSilver, err = session.GuildRoleCreate(GuildID, &discordgo.RoleParams{Name: "Silver", Color: &color, Mentionable: &mentionnable, Hoist: &hoist})
		if err != nil {
			log.Fatalf("failed to create role RoleSilver")
		}
	}
	if RoleBronze == nil {
		color := 15105570
		RoleBronze, err = session.GuildRoleCreate(GuildID, &discordgo.RoleParams{Name: "Bronze", Color: &color, Mentionable: &mentionnable, Hoist: &hoist})
		if err != nil {
			log.Fatalf("failed to create role RoleBronze")
		}
	}
	if RoleRookie == nil {
		RoleRookie, err = session.GuildRoleCreate(GuildID, &discordgo.RoleParams{Name: "Rookie", Mentionable: &mentionnable, Hoist: &hoist})
		if err != nil {
			log.Fatalf("failed to create role RoleRookie")
		}
	}
	RankRoles = []*discordgo.Role{RoleOmega, RoleChallenger, RoleDiamond, RolePlatinum, RoleGold, RoleSilver, RoleBronze}
	return err
}
