package models

type Team struct {
	Players []*Player
	OwnerID string `db:"ownerplayerID"`
	Name    string `db:"name"`
}
