package demodata

import (
	"context"
	"fmt"

	"github.com/euscs/euscs-bot/internal/db"
)

func InitWithBasicInformation() error {
	ctx := context.TODO()
	for i := 0; i < 16; i++ {
		pid1 := 12340 + i*3
		pid2 := 12340 + i*3 + 1
		pid3 := 12340 + i*3 + 2
		p1, _ := db.CreatePlayerWithID(ctx, fmt.Sprintf("%d", pid1))
		p2, _ := db.CreatePlayerWithID(ctx, fmt.Sprintf("%d", pid2))
		p3, _ := db.CreatePlayerWithID(ctx, fmt.Sprintf("%d", pid3))
		p1.CreateTeamWithName(ctx, fmt.Sprintf("team%d", i+1))
		inv2, _ := p1.Invite(ctx, p2)
		inv3, _ := p1.Invite(ctx, p3)
		inv2.Accept(ctx)
		inv3.Accept(ctx)
	}
	return nil
}
