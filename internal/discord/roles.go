package discord

import (
	"errors"

	"github.com/bwmarrin/discordgo"
)

var RoleOmega *discordgo.Role
var RoleChallenger *discordgo.Role
var RoleDiamond *discordgo.Role
var RolePlatinum *discordgo.Role
var RoleGold *discordgo.Role
var RoleSilver *discordgo.Role
var RoleBronze *discordgo.Role

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
	}
	if RoleOmega == nil {
		return errors.New("missing discord role omega in the guild")
	}
	if RoleChallenger == nil {
		return errors.New("missing discord role challenger in the guild")
	}
	if RoleDiamond == nil {
		return errors.New("missing discord role diamond in the guild")
	}
	if RolePlatinum == nil {
		return errors.New("missing discord role platinum in the guild")
	}
	if RoleGold == nil {
		return errors.New("missing discord role gold in the guild")
	}
	if RoleSilver == nil {
		return errors.New("missing discord role silver in the guild")
	}
	if RoleBronze == nil {
		return errors.New("missing discord role bronze in the guild")
	}
	RankRoles = []*discordgo.Role{RoleOmega, RoleChallenger, RoleDiamond, RolePlatinum, RoleGold, RoleSilver, RoleBronze}
	return err
}
