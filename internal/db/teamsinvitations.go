package db

import (
	"context"
	"errors"
	"time"

	"github.com/euscs/euscs-bot/internal/static"
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
	invite := TeamInvitation{Player: player, Team: team, MessageID: inviteTmp.MessageID}
	return &invite, nil
}

func (i *TeamInvitation) Save(ctx context.Context) error {
	if i.Player == nil {
		return static.ErrPlayerRequired
	}
	if i.Team == nil {
		return static.ErrTeamRequired
	}
	if i.MessageID == "" {
		return static.ErrMessageIDRequired
	}
	if i.Timestamp == 0 {
		i.Timestamp = int(time.Now().Unix())
	}
	_, err := GetTeamInvitationByID(ctx, i.MessageID)
	if err != nil && !errors.Is(err, static.ErrNotFound) {
		return err
	} else if errors.Is(err, static.ErrNotFound) {
		_, err := db.Exec("INSERT INTO teamsinvitations (playerID,team,messageID,timestamp,state) VALUES (?,?,?,?,?)", i.Player.DiscordID, i.Team.Name, i.MessageID, i.Timestamp, i.State)
		if err != nil {
			return static.ErrDB(err)
		}
	} else {
		_, err := db.NamedExec("UPDATE teamsinvitations SET state=:state WHERE messageID=:messageID", i)
		if err != nil {
			return static.ErrDB(err)
		}
	}
	return nil
}
