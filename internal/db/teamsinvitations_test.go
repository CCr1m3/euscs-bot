package db

import (
	"context"
	"errors"
	"testing"

	"github.com/euscs/euscs-bot/internal/static"
	"github.com/google/go-cmp/cmp"
)

func TestGetTeamInvitationByID(t *testing.T) {
	Clear()
	Init()
	ctx := context.TODO()
	t.Run("empty", func(t *testing.T) {
		_, err := GetTeamInvitationByID(ctx, "thisisanidforinvitation")
		if !errors.Is(err, static.ErrNotFound) {
			t.Errorf("unexpected error, should be: %s", static.ErrNotFound)
		}
	})
	p1, _ := CreatePlayerWithID(ctx, "12345")
	p2, _ := CreatePlayerWithID(ctx, "12346")
	p1.CreateTeamWithName(ctx, "team1")
	invite, err := p1.Invite(ctx, p2)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	t.Run("success", func(t *testing.T) {
		invite2, err := GetTeamInvitationByID(ctx, invite.MessageID)
		if err != nil {
			t.Errorf("unexpected error, %s", err.Error())
		}
		if !cmp.Equal(invite, invite2) {
			t.Errorf("invitations are different: %s", cmp.Diff(invite, invite2))
		}
	})
}

func TestPlayer_Invite(t *testing.T) {
	Clear()
	Init()
	ctx := context.TODO()
	p1, _ := CreatePlayerWithID(ctx, "12345")
	p2, _ := CreatePlayerWithID(ctx, "12346")
	p3, _ := CreatePlayerWithID(ctx, "12347")
	p4, _ := CreatePlayerWithID(ctx, "12348")
	t.Run("noteam", func(t *testing.T) {
		_, err := p1.Invite(ctx, p2)
		if !errors.Is(err, static.ErrNoTeam) {
			t.Errorf("unexpected error, should be: %s", static.ErrNoTeam)
		}
	})
	team1, _ := p1.CreateTeamWithName(ctx, "team1")
	t.Run("success", func(t *testing.T) {
		_, err := p1.Invite(ctx, p2)
		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}
	})
	team1.AddPlayer(ctx, p2)
	team1.AddPlayer(ctx, p3)
	t.Run("full", func(t *testing.T) {
		_, err := p1.Invite(ctx, p4)
		if !errors.Is(err, static.ErrTeamFull) {
			t.Errorf("unexpected error, should be: %s", static.ErrTeamFull)
		}
	})
	p4.CreateTeamWithName(ctx, "team4")
	t.Run("alreadyinteam", func(t *testing.T) {
		_, err := p4.Invite(ctx, p1)
		if !errors.Is(err, static.ErrUserAlreadyInTeam) {
			t.Errorf("unexpected error, should be: %s", static.ErrUserAlreadyInTeam)
		}
	})
}

func TestTeamInvitation_Accept(t *testing.T) {
	Clear()
	Init()
	ctx := context.TODO()
	p1, _ := CreatePlayerWithID(ctx, "12345")
	p2, _ := CreatePlayerWithID(ctx, "12346")
	p3, _ := CreatePlayerWithID(ctx, "12347")
	p4, _ := CreatePlayerWithID(ctx, "12348")
	p1.CreateTeamWithName(ctx, "team1")
	invite2, _ := p1.Invite(ctx, p2)
	invite3, _ := p1.Invite(ctx, p3)
	invite4, _ := p1.Invite(ctx, p4)
	t.Run("simple", func(t *testing.T) {
		invite2, _ := GetTeamInvitationByID(ctx, invite2.MessageID)
		err := invite2.Accept(ctx)
		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}
	})
	t.Run("again", func(t *testing.T) {
		invite3, _ := GetTeamInvitationByID(ctx, invite3.MessageID)
		err := invite3.Accept(ctx)
		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}
	})
	t.Run("full", func(t *testing.T) {
		invite4, _ := GetTeamInvitationByID(ctx, invite4.MessageID)
		err := invite4.Accept(ctx)
		if !errors.Is(err, static.ErrTeamFull) {
			t.Errorf("unexpected error, should be :%s", static.ErrTeamFull)
		}
	})
}

func TestTeamInvitation_Refuse(t *testing.T) {
	Clear()
	Init()
	ctx := context.TODO()
	p1, _ := CreatePlayerWithID(ctx, "12345")
	p2, _ := CreatePlayerWithID(ctx, "12346")
	p1.CreateTeamWithName(ctx, "team1")
	invite2, _ := p1.Invite(ctx, p2)
	t.Run("simple", func(t *testing.T) {
		invite2, _ := GetTeamInvitationByID(ctx, invite2.MessageID)
		err := invite2.Refuse(ctx)
		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}
	})
}
