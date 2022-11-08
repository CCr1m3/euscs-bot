package db

import (
	"context"
	"testing"
)

func Test_db_UpdatePlayer(t *testing.T) {
	clearDB()
	Init()
	ctx := context.TODO()
	err := CreatePlayer(ctx, "12345")
	if err != nil {
		t.Errorf("failed to create player: " + err.Error())
	}
	p, err := GetPlayerById(ctx, "12345")
	if err != nil {
		t.Errorf("failed to get player: " + err.Error())
	}
	p.OSUser = "osuser"
	err = UpdatePlayer(ctx, p)
	if err != nil {
		t.Errorf("failed to update player: " + err.Error())
	}
	p, err = GetPlayerByUsername(ctx, "osuser")
	if err != nil {
		t.Errorf("failed to get player: " + err.Error())
	}
	if p.OSUser != "osuser" {
		t.Errorf("failed to update player osuser: %v != osuser", p.OSUser)
	}
}
