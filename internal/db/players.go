package db

type Player struct {
	DiscordID string `db:"discordID"`
	Elo       int    `db:"elo"`
	Queuing   int    `db:"queuing"`
	Role      string `db:"role"`
	OSUser    string `db:"osuser"`
}

func CreatePlayer(discordID string) error {
	_, err := db.Exec("INSERT INTO players (discordID) VALUES (?)", discordID)
	return err
}

func GetPlayer(discordID string) (*Player, error) {
	var player Player
	err := db.Get(&player, "SELECT * FROM players WHERE discordID=?", discordID)
	if err != nil {
		return nil, err
	}
	return &player, nil
}

func (p *Player) LinkOSUser(osuser string) error {
	p.OSUser = osuser
	_, err := db.NamedExec("UPDATE players SET osuser=:osuser WHERE discordID=:discordID", p)
	return err
}

func (p *Player) AddToQueue(role string) error {
	p.Queuing = 1
	p.Role = role
	_, err := db.NamedExec("UPDATE players SET queuing=:queuing, role=:role WHERE discordID=:discordID", p)
	return err
}

func (p *Player) LeaveQueue() error {
	p.Queuing = 0
	p.Role = ""
	_, err := db.NamedExec("UPDATE players SET queuing=:queuing, role=:role WHERE discordID=:discordID", p)
	return err
}

func (p *Player) IsInQueue() bool {
	return p.Queuing == 1
}

func (p *Player) IsInMatch() (bool, error) {
	var count int
	row := db.QueryRow("SELECT COUNT(*) FROM matches JOIN matchesplayers ON matches.matchID = matchesplayers.matchID WHERE playerID=? and running=1", p.DiscordID)
	err := row.Scan(&count)
	return count > 0, err
}

func GetPlayersInQueue() ([]*Player, error) {
	players := []*Player{}
	err := db.Select(&players, "SELECT * FROM players WHERE queuing=1")
	return players, err
}

func GetGoaliesCountInQueue() (int, error) {
	var count int
	row := db.QueryRow("SELECT COUNT(discordID) FROM players WHERE queuing=1 AND (role='goalie' OR role='flex')")
	err := row.Scan(&count)
	return count, err
}

func GetForwardsCountInQueue() (int, error) {
	var count int
	row := db.QueryRow("SELECT COUNT(discordID) FROM players WHERE queuing=1 AND (role='forward' OR role='flex')")
	err := row.Scan(&count)
	return count, err
}
