package db

import (
	"context"
	"errors"
	"testing"

	"github.com/euscs/euscs-bot/internal/static"
	"github.com/google/go-cmp/cmp"
)

func TestGetPlayerByUsername(t *testing.T) {
	Clear()
	Init()
	ctx := context.TODO()
	t.Run("empty", func(t *testing.T) {
		_, err := GetPlayerByUsername(ctx, "osuser")
		if !errors.Is(err, static.ErrNotFound) {
			t.Errorf("found unexisting player")
		}
	})
	p1, _ := CreatePlayerWithID(ctx, "12345")
	p1.SetOSUser(ctx, "osuser")
	t.Run("success", func(t *testing.T) {
		p2, err := GetPlayerByUsername(ctx, "osuser")
		if err != nil {
			t.Errorf("unexpected error, %s", err.Error())
		}
		if !cmp.Equal(p1, p2) {
			t.Errorf("players are different: %s", cmp.Diff(p1, p2))
		}
	})
}

func TestGetPlayerByID(t *testing.T) {
	Clear()
	Init()
	ctx := context.TODO()
	t.Run("empty", func(t *testing.T) {
		_, err := GetPlayerByID(ctx, "12345")
		if !errors.Is(err, static.ErrNotFound) {
			t.Errorf("found unexisting player")
		}
	})
	p1, _ := CreatePlayerWithID(ctx, "12345")
	t.Run("success", func(t *testing.T) {
		p2, err := GetPlayerByID(ctx, "12345")
		if err != nil {
			t.Errorf("unexpected error, %s", err.Error())
		}
		if !cmp.Equal(p1, p2) {
			t.Errorf("players are different: %s", cmp.Diff(p1, p2))
		}
	})
}

func TestGetOrCreatePlayerByID(t *testing.T) {
	Clear()
	Init()
	ctx := context.TODO()
	t.Run("firstgetorcreate", func(t *testing.T) {
		_, err := GetOrCreatePlayerByID(ctx, "12345")
		if err != nil {
			t.Errorf("unexpected error, %s", err.Error())
		}
		_, err = GetPlayerByID(ctx, "12345")
		if err != nil {
			t.Errorf("player should be created")
		}
	})

	t.Run("secondgetorcreate", func(t *testing.T) {
		_, err := GetOrCreatePlayerByID(ctx, "12345")
		if err != nil {
			t.Errorf("unexpected error, %s", err.Error())
		}
	})
}

func TestCreatePlayerWithID(t *testing.T) {
	Clear()
	Init()
	ctx := context.TODO()
	t.Run("firstcreate", func(t *testing.T) {
		_, err := CreatePlayerWithID(ctx, "12345")
		if err != nil {
			t.Errorf("unexpected error, %s", err.Error())
		}
		_, err = GetPlayerByID(ctx, "12345")
		if err != nil {
			t.Errorf("player should be created")
		}
	})

	t.Run("secondcreate", func(t *testing.T) {
		_, err := CreatePlayerWithID(ctx, "12345")
		if !errors.Is(err, static.ErrAlreadyExists) {
			t.Errorf("unexpected error, should be: %s", static.ErrAlreadyExists)
		}
	})
}

func TestPlayer_SetElo(t *testing.T) {
	Clear()
	Init()
	ctx := context.TODO()
	t.Run("empty", func(t *testing.T) {
		p := Player{DiscordID: "12345"}
		err := p.SetElo(ctx, 1800)
		if !errors.Is(err, static.ErrNotFound) {
			t.Errorf("found unexisting player")
		}
	})
	p, _ := CreatePlayerWithID(ctx, "12345")
	t.Run("success", func(t *testing.T) {
		p.SetElo(ctx, 1800)
		p2, _ := GetPlayerByID(ctx, "12345")
		if !cmp.Equal(p, p2) {
			t.Errorf("players are different: %s", cmp.Diff(p, p2))
		}
	})
}

func TestPlayer_SetOSUser(t *testing.T) {
	Clear()
	Init()
	ctx := context.TODO()
	t.Run("empty", func(t *testing.T) {
		p := Player{DiscordID: "12345"}
		err := p.SetOSUser(ctx, "OSUser")
		if !errors.Is(err, static.ErrNotFound) {
			t.Errorf("found unexisting player")
		}
	})
	p, _ := CreatePlayerWithID(ctx, "12345")
	t.Run("success", func(t *testing.T) {
		p.SetOSUser(ctx, "OSUser")
		p2, _ := GetPlayerByID(ctx, "12345")
		if !cmp.Equal(p, p2) {
			t.Errorf("players are different: %s", cmp.Diff(p, p2))
		}
	})
}
