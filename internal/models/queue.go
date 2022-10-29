package models

type Role string

const (
	RoleFlex    Role = "flex"
	RoleForward Role = "forward"
	RoleGoalie  Role = "goalie"
)

type QueuedPlayer struct {
	Player
	Role      Role
	EntryTime int
}
