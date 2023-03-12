package db

import (
	"context"
	"errors"
	"testing"

	"github.com/euscs/euscs-bot/internal/static"
	"github.com/google/go-cmp/cmp"
)

func TestGetPlayerByUsername(t *testing.T) {
	clearDB()
	Init()
	ctx := context.TODO()
	t.Run("empty", func(t *testing.T) {
		_, err := GetPlayerByUsername(ctx, "osuser")
		if !errors.Is(err, static.ErrNotFound) {
			t.Errorf("found unexisting player")
		}
	})
	p1 := &Player{DiscordID: "12345", OSUser: "osuser"}
	p1.Save(ctx)
	t.Run("success", func(t *testing.T) {
		p2, err := GetPlayerByUsername(ctx, "osuser")
		if err != nil {
			t.Errorf("unexpected error, %s", err.Error())
		}
		if !cmp.Equal(p1, p2) {
			t.Logf("want: %#v\n", p1)
			t.Logf("got: %#v\n", p2)
			t.Errorf("players are different")
		}
	})
}

func TestGetPlayerByID(t *testing.T) {
	clearDB()
	Init()
	ctx := context.TODO()
	t.Run("empty", func(t *testing.T) {
		_, err := GetPlayerByID(ctx, "12345")
		if !errors.Is(err, static.ErrNotFound) {
			t.Errorf("found unexisting player")
		}
	})
	p1 := &Player{DiscordID: "12345"}
	p1.Save(ctx)
	t.Run("success", func(t *testing.T) {
		p2, err := GetPlayerByID(ctx, "12345")
		if err != nil {
			t.Errorf("unexpected error, %s", err.Error())
		}
		if !cmp.Equal(p1, p2) {
			t.Logf("want: %#v\n", p1)
			t.Logf("got: %#v\n", p2)
			t.Errorf("players are different")
		}
	})
}

func TestPlayer_Save(t *testing.T) {
	clearDB()
	Init()
	ctx := context.TODO()
	t.Run("simplesave", func(t *testing.T) {
		p := &Player{DiscordID: "12345"}
		err := p.Save(ctx)
		if err != nil {
			t.Errorf("unexpected error, %s", err.Error())
		}
	})
	t.Run("savewithdiscourdid", func(t *testing.T) {
		p := &Player{}
		err := p.Save(ctx)
		if !errors.Is(err, static.ErrDiscordIDRequired) {
			t.Errorf("error should be: %s", static.ErrDiscordIDRequired)
		}
	})
	t.Run("saveandedit", func(t *testing.T) {
		p := &Player{DiscordID: "12346"}
		err := p.Save(ctx)
		if err != nil {
			t.Errorf("unexpected error, %s", err.Error())
		}
		p.Elo = 1600
		p.TwitchID = "1233456"
		p.OSUser = "osuser"
		err = p.Save(ctx)
		if err != nil {
			t.Errorf("unexpected error, %s", err.Error())
		}
	})
}

func TestGetOrCreatePlayerByID(t *testing.T) {
	clearDB()
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
