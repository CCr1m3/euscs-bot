package models

type Player struct {
	DiscordID string `db:"discordID"`
	TwitchID  string `db:"twitchID"`
	Elo       int    `db:"elo"`
	OSUser    string `db:"osuser"`
}
