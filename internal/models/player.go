package models

type Player struct {
	DiscordID      string `db:"discordID"`
	Elo            int    `db:"elo"`
	OSUser         string `db:"osuser"`
	LastRankUpdate int    `db:"lastrankupdate"`
	Credits        int    `db:"credits"`
}
