package db

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/euscs/euscs-bot/internal/discord"
	"github.com/euscs/euscs-bot/internal/static"
	"github.com/google/uuid"
)

type InvitationState int

const (
	InvitationPending  = 0
	InvitationRefused  = 1
	InvitationAccepted = 2
)

type TeamInvitation struct {
	Player    *Player
	Team      *Team
	MessageID string          `db:"messageID"`
	Timestamp int             `db:"timestamp"`
	State     InvitationState `db:"state"`
}

func GetTeamInvitationByID(ctx context.Context, messageID string) (*TeamInvitation, error) {
	var inviteTmp struct {
		TeamInvitation
		PlayerID string `db:"playerID"`
		Team     string `db:"team"`
	}
	err := db.Get(&inviteTmp, "SELECT playerID,team,messageID,timestamp,state FROM teamsinvitations WHERE messageID=?", messageID)
	if err != nil {
		return nil, static.ErrDB(err)
	}
	player, err := GetPlayerByID(ctx, inviteTmp.PlayerID)
	if err != nil {
		return nil, err
	}
	team, err := GetTeamByName(ctx, inviteTmp.Team)
	if err != nil {
		return nil, err
	}
	invite := TeamInvitation{Player: player, Team: team, MessageID: inviteTmp.MessageID, Timestamp: inviteTmp.Timestamp}
	return &invite, nil
}

func (p *Player) Invite(ctx context.Context, p2 *Player) (*TeamInvitation, error) {
	team, err := p.GetTeam(ctx)
	if err != nil {
		if errors.Is(err, static.ErrNotFound) {
			return nil, static.ErrNoTeam
		}
		return nil, err
	}
	if len(team.Players) >= 3 {
		return nil, static.ErrTeamFull
	}
	team2, err := p2.GetTeam(ctx)
	if err != nil && !errors.Is(err, static.ErrNotFound) {
		return nil, err
	} else if team2 != nil {
		return nil, static.ErrUserAlreadyInTeam
	}

	var messageID string
	if p2.isDummy() {
		messageID = uuid.New().String()
	} else {
		session := discord.GetSession()
		channel, err := session.UserChannelCreate(p2.DiscordID)
		if err != nil {
			return nil, err
		}
		messageContent := fmt.Sprintf("You have been invited by %s to the team '%s'", "<@"+p.DiscordID+">", team.Name)
		message, err := session.ChannelMessageSendComplex(channel.ID, &discordgo.MessageSend{
			Content: messageContent,
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							Label:    "Accept",
							Style:    discordgo.PrimaryButton,
							CustomID: "accept_invite",
						},
						discordgo.Button{
							Label:    "Refuse",
							Style:    discordgo.DangerButton,
							CustomID: "refuse_invite",
						},
					},
				},
			}})
		if err != nil {
			return nil, err
		}
		messageID = message.ID
	}
	timestamp := int(time.Now().Unix())
	_, err = db.Exec("INSERT INTO teamsinvitations (playerID,team,messageID,timestamp,state) VALUES (?,?,?,?,?)", p2.DiscordID, team.Name, messageID, timestamp, InvitationPending)
	if err != nil {
		return nil, static.ErrDB(err)
	}
	return &TeamInvitation{Player: p2, Team: team, MessageID: messageID, Timestamp: timestamp, State: InvitationPending}, nil
}

func (ti *TeamInvitation) Accept(ctx context.Context) error {
	return ti.Team.AddPlayer(ctx, ti.Player)
}

func (ti *TeamInvitation) Refuse(ctx context.Context) error {
	return nil
}
