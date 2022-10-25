package rank

import (
	"database/sql"
	"errors"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/haashi/omega-strikers-bot/internal/db"
	"github.com/haashi/omega-strikers-bot/internal/discord"
	"github.com/haashi/omega-strikers-bot/internal/models"
	log "github.com/sirupsen/logrus"
)

func Init() {
	log.Info("starting rank service")
}

func getOrCreatePlayer(playerID string) (*models.Player, error) {
	p, err := db.GetPlayerById(playerID)
	if err != nil {
		err = db.CreatePlayer(playerID)
		if err != nil {
			return nil, err
		}
		p, err = db.GetPlayerById(playerID)
		if err != nil {
			return nil, err
		}
	}
	return p, nil
}

func LinkPlayerToUsername(playerID string, username string) error {
	player, err := getOrCreatePlayer(playerID)
	if err != nil {
		return err
	}
	if player.OSUser == "" {
		_, err := db.GetPlayerByUsername(username)
		if !errors.Is(err, sql.ErrNoRows) {
			return &models.UsernameAlreadyLinkedError{Username: username}
		}
		player.OSUser = username
		err = db.UpdatePlayer(player)
		if err != nil {
			return err
		}
		err = UpdateRank(playerID, true)
		return err
	} else {
		return &models.UserAlreadyLinkedError{UserID: playerID}
	}
}

func UnlinkPlayer(playerID string) error {
	player, err := getOrCreatePlayer(playerID)
	if err != nil {
		return err
	}
	if player.OSUser == "" {
		return &models.NotLinkedError{UserID: playerID}
	}
	player.LastRankUpdate = 0
	player.OSUser = ""
	err = db.UpdatePlayer(player)
	if err != nil {
		return err
	}
	return err
}

func GetLinkedUsername(playerID string) (string, error) {
	player, err := getOrCreatePlayer(playerID)
	if err != nil {
		return "", err
	}
	return player.OSUser, nil
}

func GetLinkedUser(username string) (string, error) {
	player, err := db.GetPlayerByUsername(username)
	if err != nil {
		return "", err
	}
	return player.DiscordID, nil
}

func UpdateRankIfNeeded(playerID string) error {
	player, err := getOrCreatePlayer(playerID)
	if err != nil {
		return err
	}
	if player.OSUser == "" {
		return &models.NotLinkedError{UserID: playerID}
	}
	updateDelay := time.Hour * 24
	if os.Getenv("mode") == "dev" {
		updateDelay = time.Hour * 1
	}
	if time.Since(time.Unix(int64(player.LastRankUpdate), 0)) > updateDelay {
		return UpdateRank(player.DiscordID, false)
	} else {
		return &models.RankUpdateTooFastError{UserID: playerID}
	}
}

func UpdateRank(playerID string, updateDiscordRole bool) error {
	player, err := getOrCreatePlayer(playerID)
	if err != nil {
		return err
	}
	if player.OSUser == "" {
		return &models.NotLinkedError{UserID: playerID}
	}
	log.Infof("updating player elo %s", player.DiscordID)
	rank, err := GetRankFromUsername(player.OSUser)
	if err != nil {
		log.Errorf("failed to retrieve rank of player %s: "+err.Error(), player.DiscordID)
		return err
	}
	player.Elo = rank
	player.LastRankUpdate = int(time.Now().Unix())
	err = db.UpdatePlayer(player)
	if err != nil {
		log.Errorf("failed to update player %s: "+err.Error(), player.DiscordID)
	}
	if updateDiscordRole {
		go func() { //update in background
			err := updatePlayerDiscordRole(player.DiscordID)
			if err != nil {
				log.Errorf("failed to update discord role of user %s: "+err.Error(), player.DiscordID)
			}
		}()
	}
	return err
}

func updatePlayerDiscordRole(playerID string) error {
	session := discord.GetSession()
	guildID := os.Getenv("guildid")
	player, err := db.GetPlayerById(playerID)
	if err != nil {
		return err
	}
	var roleToAdd *discordgo.Role
	if player.Elo >= 2900 {
		roleToAdd = discord.RoleOmega
	} else if player.Elo >= 2600 {
		roleToAdd = discord.RoleChallenger
	} else if player.Elo >= 2300 {
		roleToAdd = discord.RoleDiamond
	} else if player.Elo >= 2000 {
		roleToAdd = discord.RolePlatinum
	} else if player.Elo >= 1700 {
		roleToAdd = discord.RoleGold
	} else if player.Elo >= 1400 {
		roleToAdd = discord.RoleSilver
	} else if player.Elo >= 1100 {
		roleToAdd = discord.RoleBronze
	} else {
		return nil
	}
	member, err := session.GuildMember(guildID, player.DiscordID)
	if err != nil {
		return err
	}
	var currentRole *discordgo.Role
	for _, roleID := range member.Roles {
		if roleID == discord.RoleOmega.ID {
			currentRole = discord.RoleOmega
		}
		if roleID == discord.RoleChallenger.ID {
			currentRole = discord.RoleChallenger
		}
		if roleID == discord.RoleDiamond.ID {
			currentRole = discord.RoleDiamond
		}
		if roleID == discord.RolePlatinum.ID {
			currentRole = discord.RolePlatinum
		}
		if roleID == discord.RoleGold.ID {
			currentRole = discord.RoleGold
		}
		if roleID == discord.RoleSilver.ID {
			currentRole = discord.RoleSilver
		}
		if roleID == discord.RoleBronze.ID {
			currentRole = discord.RoleBronze
		}
	}
	if currentRole != nil && currentRole.Position > roleToAdd.Position {
		//we only update for peak elo
		return nil
	}
	for _, rankRole := range discord.RankRoles {
		err := session.GuildMemberRoleRemove(guildID, player.DiscordID, rankRole.ID)
		if err != nil {
			return err
		}
	}
	err = session.GuildMemberRoleAdd(guildID, player.DiscordID, roleToAdd.ID)
	return err
}
