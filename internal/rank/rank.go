package rank

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/euscs/euscs-bot/internal/db"
	"github.com/euscs/euscs-bot/internal/discord"
	"github.com/euscs/euscs-bot/internal/static"
	log "github.com/sirupsen/logrus"
)

func LinkPlayerToUsername(ctx context.Context, playerID string, username string) error {
	player, err := db.GetOrCreatePlayerByID(ctx, playerID)
	if err != nil {
		return err
	}
	if player.OSUser == "" {
		_, err := db.GetPlayerByUsername(ctx, username)
		if err == nil {
			return static.ErrUsernameAlreadyLinked
		} else if err != nil && !errors.Is(err, static.ErrNotFound) {
			return err
		}
		err = player.SetOSUser(ctx, username)
		if err != nil {
			return err
		}
		err = UpdateRank(ctx, playerID)
		if err != nil {
			return err
		}
		return nil
	} else {
		return static.ErrUserAlreadyLinked
	}
}

func UnlinkPlayer(ctx context.Context, playerID string) error {
	player, err := db.GetOrCreatePlayerByID(ctx, playerID)
	if err != nil {
		return err
	}
	if player.OSUser == "" {
		return static.ErrUserNotLinked
	}
	err = player.SetElo(ctx, 0)
	if err != nil {
		return err
	}
	err = player.SetOSUser(ctx, "")
	if err != nil {
		return err
	}
	go func() { //update in background
		err := updatePlayerDiscordRole(ctx, player.DiscordID)
		if err != nil {
			log.Errorf("failed to update discord role of user %s: "+err.Error(), player.DiscordID)
		}
	}()
	return err
}

func GetLinkedUsername(ctx context.Context, playerID string) (string, error) {
	player, err := db.GetOrCreatePlayerByID(ctx, playerID)
	if err != nil {
		return "", err
	}
	return player.OSUser, nil
}

func GetLinkedUser(ctx context.Context, username string) (string, error) {
	player, err := db.GetPlayerByUsername(ctx, username)
	if err != nil {
		return "", err
	}
	return player.DiscordID, nil
}

func UpdateRankIfNeeded(ctx context.Context, playerID string) error {
	player, err := db.GetOrCreatePlayerByID(ctx, playerID)
	if err != nil {
		return err
	}
	if player.OSUser == "" {
		return static.ErrUserNotLinked
	}
	updateDelay := time.Hour * 1
	if os.Getenv("mode") == "dev" {
		updateDelay = time.Second * 15
	}
	if time.Since(time.Unix(int64(player.LastRankUpdate), 0)) > updateDelay {
		return UpdateRank(ctx, player.DiscordID)
	} else {
		return static.ErrRankUpdateTooFast
	}
}

func UpdateRank(ctx context.Context, playerID string) error {
	player, err := db.GetOrCreatePlayerByID(ctx, playerID)
	if err != nil {
		return err
	}
	if player.OSUser == "" {
		return static.ErrUserNotLinked
	}
	log.Infof("updating player elo %s", player.DiscordID)
	info, err := GetCorestrikeInfoFromUsername(ctx, player.OSUser)
	if err != nil {
		log.Errorf("failed to retrieve rank of player %s: "+err.Error(), player.DiscordID)
		if errors.Is(err, static.ErrUsernameInvalid) {
			log.Warningf("unlinking %s because username %s was not valid", playerID, player.OSUser)
			if player.OSUser != "" {
				err2 := UnlinkPlayer(ctx, playerID)
				if err2 != nil {
					log.Errorf("failed to unlink player %s: "+err.Error(), playerID)
				}
			}
			return err
		}
		return err
	}
	rank := info.RankedStats.Rating
	if rank > player.Elo {
		player.SetElo(ctx, rank)
		if err != nil {
			log.Errorf("failed to update player %s: "+err.Error(), player.DiscordID)
		}
	}

	go func() { //update in background
		err := updatePlayerDiscordRole(ctx, player.DiscordID)
		if err != nil {
			log.Errorf("failed to update discord role of user %s: "+err.Error(), player.DiscordID)
		}
	}()
	return err
}

func updatePlayerDiscordRole(ctx context.Context, playerID string) error {
	session := discord.GetSession()
	guildID := discord.GuildID
	player, err := db.GetPlayerByID(ctx, playerID)
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
	} else if player.Elo > 0 {
		roleToAdd = discord.RoleRookie
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
		if roleID == discord.RoleRookie.ID {
			currentRole = discord.RoleRookie
		}
	}
	if currentRole != nil && roleToAdd != nil && currentRole.Position >= roleToAdd.Position {
		//we only update for peak elo
		return nil
	}
	for _, rankRole := range discord.RankRoles {
		err := session.GuildMemberRoleRemove(guildID, player.DiscordID, rankRole.ID)
		if err != nil {
			return err
		}
	}
	if roleToAdd != nil {
		err = session.GuildMemberRoleAdd(guildID, player.DiscordID, roleToAdd.ID)
		if err != nil {
			return err
		}
		_, err = session.ChannelMessageSend(discord.RankUpChannel.ID, fmt.Sprintf("%s just got promoted to %s !", "<@"+player.DiscordID+">", roleToAdd.Name))
		if err != nil {
			return err
		}
	}
	return err
}

func Init() {

}
