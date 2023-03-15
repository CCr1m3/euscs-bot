package slashcommands

import (
	"context"
	"errors"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/euscs/euscs-bot/internal/db"
	"github.com/euscs/euscs-bot/internal/static"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type TeamCreate struct{}

func (p TeamCreate) Name() string {
	return "teamcreate"
}

func (p TeamCreate) Description() string {
	return "Allow you to create a team."
}

func (p TeamCreate) RequiredPerm() *int64 {
	perm := int64(discordgo.PermissionSendMessages)
	return &perm
}

func (p TeamCreate) Options() []*discordgo.ApplicationCommandOption {
	return []*discordgo.ApplicationCommandOption{
		{
			Name:        "teamname",
			Description: "name of the team",
			Type:        discordgo.ApplicationCommandOptionString,
			Required:    true,
		},
	}
}
func (p TeamCreate) Run(s *discordgo.Session, i *discordgo.InteractionCreate) {
	ctx := context.WithValue(context.Background(), static.UUIDKey, uuid.New())
	options := i.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}
	var teamName string
	if val, ok := optionMap["teamname"]; ok {
		teamName = val.StringValue()
	}
	log.WithFields(log.Fields{
		string(static.UUIDKey):     ctx.Value(static.UUIDKey),
		string(static.CallerIDKey): i.Member.User.ID,
		string(static.TeamNameKey): teamName,
	}).Info("teamcreate slash command invoked")
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "teamcreate slash command invoked. Please wait...",
		},
	})
	if err != nil {
		log.WithFields(log.Fields{
			string(static.UUIDKey):  ctx.Value(static.UUIDKey),
			string(static.ErrorKey): err.Error(),
		}).Error("failed to send message")
		return
	}
	var message string
	defer func() {
		_, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: &message,
		})
		if err != nil {
			log.WithFields(log.Fields{
				string(static.UUIDKey):  ctx.Value(static.UUIDKey),
				string(static.ErrorKey): err.Error(),
			}).Error("failed to edit message")
		}
	}()
	if teamName == "" {
		message = "Please enter a teamname."
		log.WithFields(log.Fields{
			string(static.UUIDKey):     ctx.Value(static.UUIDKey),
			string(static.CallerIDKey): i.Member.User.ID,
		}).Warning("teamcreate failed, no arguments")
		return
	}
	player, err := db.GetOrCreatePlayerByID(ctx, i.Member.User.ID)
	if err != nil {
		log.WithFields(log.Fields{
			string(static.UUIDKey):     ctx.Value(static.UUIDKey),
			string(static.CallerIDKey): i.Member.User.ID,
			string(static.ErrorKey):    err.Error(),
		}).Error("failed to get player")
		message = fmt.Sprintf("Failed to created team: %s", teamName)
	}
	_, err = player.CreateTeamWithName(ctx, teamName)
	if err == nil {
		message = fmt.Sprintf("Succesfully created team: %s", teamName)
	} else {
		log.WithFields(log.Fields{
			string(static.UUIDKey):     ctx.Value(static.UUIDKey),
			string(static.CallerIDKey): i.Member.User.ID,
			string(static.ErrorKey):    err.Error(),
			string(static.TeamNameKey): teamName,
		}).Warning("failed to created team")
		message = fmt.Sprintf("Failed to created team: %s.", teamName)
		switch {
		case errors.Is(err, static.ErrTeamnameTaken):
			message += " A team already exists with this name."
		case errors.Is(err, static.ErrUserAlreadyInTeam):
			message += " You already have a team."
		default:
			message += " Unexpected Error."
		}
	}
}

type TeamInvite struct{}

func (p TeamInvite) Name() string {
	return "teaminvite"
}

func (p TeamInvite) Description() string {
	return "Allow you to invite an user to your team."
}

func (p TeamInvite) RequiredPerm() *int64 {
	perm := int64(discordgo.PermissionSendMessages)
	return &perm
}

func (p TeamInvite) Options() []*discordgo.ApplicationCommandOption {
	return []*discordgo.ApplicationCommandOption{
		{
			Name:        "discorduser",
			Description: "User in discord",
			Type:        discordgo.ApplicationCommandOptionUser,
			Required:    true,
		},
	}
}
func (p TeamInvite) Run(s *discordgo.Session, i *discordgo.InteractionCreate) {
	ctx := context.WithValue(context.Background(), static.UUIDKey, uuid.New())
	options := i.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}
	var invitedPlayerID string
	if val, ok := optionMap["discorduser"]; ok {
		invitedPlayerID = val.UserValue(s).ID
	}
	log.WithFields(log.Fields{
		string(static.UUIDKey):     ctx.Value(static.UUIDKey),
		string(static.CallerIDKey): i.Member.User.ID,
		string(static.PlayerIDKey): invitedPlayerID,
	}).Info("invite slash command invoked")
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Invite slash command invoked. Please wait...",
		},
	})
	if err != nil {
		log.WithFields(log.Fields{
			string(static.UUIDKey):  ctx.Value(static.UUIDKey),
			string(static.ErrorKey): err.Error(),
		}).Error("failed to send message")
		return
	}
	var message string
	defer func() {
		_, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: &message,
		})
		if err != nil {
			log.WithFields(log.Fields{
				string(static.UUIDKey):  ctx.Value(static.UUIDKey),
				string(static.ErrorKey): err.Error(),
			}).Error("failed to edit message")
		}
	}()
	if invitedPlayerID == "" {
		message = "Please enter user in discord."
		log.WithFields(log.Fields{
			string(static.UUIDKey):     ctx.Value(static.UUIDKey),
			string(static.CallerIDKey): i.Member.User.ID,
		}).Warning("invite failed, no arguments")
		return
	}
	player, err := db.GetOrCreatePlayerByID(ctx, i.Member.User.ID)
	if err != nil {
		log.WithFields(log.Fields{
			string(static.UUIDKey):     ctx.Value(static.UUIDKey),
			string(static.CallerIDKey): i.Member.User.ID,
			string(static.ErrorKey):    err.Error(),
		}).Error("failed to get player")
		message = "Failed to invite."
	}
	player2, err := db.GetOrCreatePlayerByID(ctx, invitedPlayerID)
	if err != nil {
		log.WithFields(log.Fields{
			string(static.UUIDKey):     ctx.Value(static.UUIDKey),
			string(static.CallerIDKey): i.Member.User.ID,
			string(static.ErrorKey):    err.Error(),
		}).Error("failed to get invited player")
		message = "Failed to invite."
	}
	_, err = player.Invite(ctx, player2)
	if err == nil {
		message = "Successfully invited."
	} else {
		log.WithFields(log.Fields{
			string(static.UUIDKey):     ctx.Value(static.UUIDKey),
			string(static.CallerIDKey): i.Member.User.ID,
			string(static.PlayerIDKey): invitedPlayerID,
			string(static.ErrorKey):    err.Error(),
		}).Warning("invite failed")
		message = "Failed to invite."
		switch {
		case errors.Is(err, static.ErrTeamFull):
			message += " Your team is full."
		case errors.Is(err, static.ErrUserAlreadyInTeam):
			message += " This user already has a team."
		case errors.Is(err, static.ErrNoTeam):
			message += " You don't have a team."
		case errors.Is(err, static.ErrNotTeamOwner):
			message += " You are not the team owner."
		default:
			log.WithFields(log.Fields{
				string(static.UUIDKey):     ctx.Value(static.UUIDKey),
				string(static.CallerIDKey): i.Member.User.ID,
				string(static.PlayerIDKey): invitedPlayerID,
				string(static.ErrorKey):    err.Error(),
			}).Error("invite failed")
			message += " Unexpected Error."
		}
	}
}

type TeamKick struct{}

func (p TeamKick) Name() string {
	return "teamkick"
}

func (p TeamKick) Description() string {
	return "Allow you to kick an user from your team."
}

func (p TeamKick) RequiredPerm() *int64 {
	perm := int64(discordgo.PermissionSendMessages)
	return &perm
}

func (p TeamKick) Options() []*discordgo.ApplicationCommandOption {
	return []*discordgo.ApplicationCommandOption{
		{
			Name:        "discorduser",
			Description: "User in discord",
			Type:        discordgo.ApplicationCommandOptionUser,
			Required:    true,
		},
	}
}

func (p TeamKick) Run(s *discordgo.Session, i *discordgo.InteractionCreate) {
	ctx := context.WithValue(context.Background(), static.UUIDKey, uuid.New())
	options := i.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}
	var kickedPlayerID string
	if val, ok := optionMap["discorduser"]; ok {
		kickedPlayerID = val.UserValue(s).ID
	}
	log.WithFields(log.Fields{
		string(static.UUIDKey):     ctx.Value(static.UUIDKey),
		string(static.CallerIDKey): i.Member.User.ID,
		string(static.PlayerIDKey): kickedPlayerID,
	}).Info("kick slash command invoked")
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "kick slash command invoked. Please wait...",
		},
	})
	if err != nil {
		log.WithFields(log.Fields{
			string(static.UUIDKey):  ctx.Value(static.UUIDKey),
			string(static.ErrorKey): err.Error(),
		}).Error("failed to send message")
		return
	}
	var message string
	defer func() {
		_, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: &message,
		})
		if err != nil {
			log.WithFields(log.Fields{
				string(static.UUIDKey):  ctx.Value(static.UUIDKey),
				string(static.ErrorKey): err.Error(),
			}).Error("failed to edit message")
		}
	}()
	if kickedPlayerID == "" {
		message = "Please enter user in discord."
		log.WithFields(log.Fields{
			string(static.UUIDKey):     ctx.Value(static.UUIDKey),
			string(static.CallerIDKey): i.Member.User.ID,
		}).Warning("kick failed, no arguments")
		return
	}
	player, err := db.GetOrCreatePlayerByID(ctx, i.Member.User.ID)
	if err != nil {
		log.WithFields(log.Fields{
			string(static.UUIDKey):     ctx.Value(static.UUIDKey),
			string(static.CallerIDKey): i.Member.User.ID,
			string(static.ErrorKey):    err.Error(),
		}).Error("failed to get player")
		message = "Failed to kick."
	}
	player2, err := db.GetOrCreatePlayerByID(ctx, kickedPlayerID)
	if err != nil {
		log.WithFields(log.Fields{
			string(static.UUIDKey):     ctx.Value(static.UUIDKey),
			string(static.CallerIDKey): i.Member.User.ID,
			string(static.ErrorKey):    err.Error(),
		}).Error("failed to get kicked player")
		message = "Failed to kick."
	}
	err = player.KickPlayerFromTeam(ctx, player2)
	if err == nil {
		message = "Successfully kicked."
	} else {
		log.WithFields(log.Fields{
			string(static.UUIDKey):     ctx.Value(static.UUIDKey),
			string(static.CallerIDKey): i.Member.User.ID,
			string(static.PlayerIDKey): kickedPlayerID,
			string(static.ErrorKey):    err.Error(),
		}).Warning("kick failed")
		message = "Failed to kick."
		switch {
		case errors.Is(err, static.ErrPlayerNotInTeam):
			message += " This user is not in your team."
		case errors.Is(err, static.ErrNoTeam):
			message += " You don't have a team."
		case errors.Is(err, static.ErrNotTeamOwner):
			message += " You are not the team owner."
		case errors.Is(err, static.ErrOwnerNotInTeam):
			message += " You can't kick yourself."
		default:
			log.WithFields(log.Fields{
				string(static.UUIDKey):     ctx.Value(static.UUIDKey),
				string(static.CallerIDKey): i.Member.User.ID,
				string(static.PlayerIDKey): kickedPlayerID,
				string(static.ErrorKey):    err.Error(),
			}).Error("kick failed")
			message += " Unexpected Error."
		}
	}
}

type TeamInfo struct{}

func (p TeamInfo) Name() string {
	return "teaminfo"
}

func (p TeamInfo) Description() string {
	return "Allow you to display info about your team."
}

func (p TeamInfo) RequiredPerm() *int64 {
	perm := int64(discordgo.PermissionSendMessages)
	return &perm
}

func (p TeamInfo) Options() []*discordgo.ApplicationCommandOption {
	return []*discordgo.ApplicationCommandOption{}
}

func (p TeamInfo) Run(s *discordgo.Session, i *discordgo.InteractionCreate) {
	ctx := context.WithValue(context.Background(), static.UUIDKey, uuid.New())
	options := i.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}
	log.WithFields(log.Fields{
		string(static.UUIDKey):     ctx.Value(static.UUIDKey),
		string(static.CallerIDKey): i.Member.User.ID,
	}).Info("team info slash command invoked")
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "team info slash command invoked. Please wait...",
		},
	})
	if err != nil {
		log.WithFields(log.Fields{
			string(static.UUIDKey):  ctx.Value(static.UUIDKey),
			string(static.ErrorKey): err.Error(),
		}).Error("failed to send message")
		return
	}
	var message string
	var embed *discordgo.MessageEmbed
	defer func() {
		var edit *discordgo.WebhookEdit
		if embed != nil {
			edit = &discordgo.WebhookEdit{
				Content: &message,
				Embeds: &[]*discordgo.MessageEmbed{
					embed,
				},
			}
		} else {
			edit = &discordgo.WebhookEdit{
				Content: &message,
			}
		}
		_, err = s.InteractionResponseEdit(i.Interaction, edit)
		if err != nil {
			log.WithFields(log.Fields{
				string(static.UUIDKey):  ctx.Value(static.UUIDKey),
				string(static.ErrorKey): err.Error(),
			}).Error("failed to edit message")
		}
	}()
	player, err := db.GetOrCreatePlayerByID(ctx, i.Member.User.ID)
	if err != nil {
		log.WithFields(log.Fields{
			string(static.UUIDKey):     ctx.Value(static.UUIDKey),
			string(static.CallerIDKey): i.Member.User.ID,
			string(static.ErrorKey):    err.Error(),
		}).Error("failed to get player")
		message = "Failed to get team info."
	}
	team, err := player.GetTeam(ctx)
	if err == nil {
		playersMsg := ""
		for _, p := range team.Players {
			line := "<@" + p.DiscordID + ">\n"
			if team.OwnerID == p.DiscordID {
				line = "👑" + line
			}
			playersMsg += line
		}
		embed = &discordgo.MessageEmbed{
			Title: team.Name,
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:  "Players:",
					Value: playersMsg,
				},
			},
		}
	} else {
		log.WithFields(log.Fields{
			string(static.UUIDKey):     ctx.Value(static.UUIDKey),
			string(static.CallerIDKey): i.Member.User.ID,
			string(static.ErrorKey):    err.Error(),
		}).Warning("team info failed")
		message = "Failed to get team info."
		switch {
		case errors.Is(err, static.ErrNotFound):
			message += " You don't have a team."
		default:
			log.WithFields(log.Fields{
				string(static.UUIDKey):     ctx.Value(static.UUIDKey),
				string(static.CallerIDKey): i.Member.User.ID,
				string(static.ErrorKey):    err.Error(),
			}).Error("team info failed")
			message += " Unexpected Error."
		}
	}
}

type TeamLeave struct{}

func (p TeamLeave) Name() string {
	return "teamleave"
}

func (p TeamLeave) Description() string {
	return "Allow you leave your team. If you are owner, the team is deleted."
}

func (p TeamLeave) RequiredPerm() *int64 {
	perm := int64(discordgo.PermissionSendMessages)
	return &perm
}

func (p TeamLeave) Options() []*discordgo.ApplicationCommandOption {
	return []*discordgo.ApplicationCommandOption{}
}

func (p TeamLeave) Run(s *discordgo.Session, i *discordgo.InteractionCreate) {
	ctx := context.WithValue(context.Background(), static.UUIDKey, uuid.New())
	options := i.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}
	log.WithFields(log.Fields{
		string(static.UUIDKey):     ctx.Value(static.UUIDKey),
		string(static.CallerIDKey): i.Member.User.ID,
	}).Info("team leave slash command invoked")
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "team leave slash command invoked. Please wait...",
		},
	})
	if err != nil {
		log.WithFields(log.Fields{
			string(static.UUIDKey):  ctx.Value(static.UUIDKey),
			string(static.ErrorKey): err.Error(),
		}).Error("failed to send message")
		return
	}
	var message string
	defer func() {
		_, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: &message,
		})
		if err != nil {
			log.WithFields(log.Fields{
				string(static.UUIDKey):  ctx.Value(static.UUIDKey),
				string(static.ErrorKey): err.Error(),
			}).Error("failed to edit message")
		}
	}()
	player, err := db.GetOrCreatePlayerByID(ctx, i.Member.User.ID)
	if err != nil {
		log.WithFields(log.Fields{
			string(static.UUIDKey):     ctx.Value(static.UUIDKey),
			string(static.CallerIDKey): i.Member.User.ID,
			string(static.ErrorKey):    err.Error(),
		}).Error("failed to get player")
		message = "Failed to leave team."
	}
	err = player.LeaveTeam(ctx)
	if err != nil {
		log.WithFields(log.Fields{
			string(static.UUIDKey):     ctx.Value(static.UUIDKey),
			string(static.CallerIDKey): i.Member.User.ID,
			string(static.ErrorKey):    err.Error(),
		}).Warning("team info failed")
		message = "Failed to get team info."
		switch {
		case errors.Is(err, static.ErrNoTeam):
			message += " You don't have a team."
		default:
			log.WithFields(log.Fields{
				string(static.UUIDKey):     ctx.Value(static.UUIDKey),
				string(static.CallerIDKey): i.Member.User.ID,
				string(static.ErrorKey):    err.Error(),
			}).Error("team info failed")
			message += " Unexpected Error."
		}
	} else {
		message = "Successfully left team."
	}
}
