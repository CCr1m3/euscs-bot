package static

import (
	"database/sql"
	"errors"
	"fmt"
)

func ErrDB(err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound
	}
	return fmt.Errorf("database error: %w", err)
}

var ErrUsernameInvalid = errors.New("username invalid")
var ErrUserNotLinked = errors.New("user not linked")
var ErrAlreadyExists = errors.New("already exists")
var ErrUserAlreadyLinked = errors.New("user is already linked")
var ErrUsernameAlreadyLinked = errors.New("username is already linked")
var ErrNotFound = errors.New("not found")
var ErrNoTeam = errors.New("no team")
var ErrPlayerNotInTeam = errors.New("not in team")
var ErrTeamFull = errors.New("team full")
var ErrUserAlreadyInTeam = errors.New("user already in team")
var ErrTeamnameTaken = errors.New("team name is taken")
var ErrNotTeamOwner = errors.New("not team owner")
var ErrOwnerNotInTeam = errors.New("owner not in team")
var ErrCorestrikeNotFound = errors.New("username not found in corestrike")
var ErrUnrankedUser = errors.New("corestrike only found unranked user")
var ErrRankUpdateTooFast = errors.New("rank update too fast")
var ErrDiscordIDRequired = errors.New("discord id required")
var ErrMessageIDRequired = errors.New("message id required")
var ErrPlayerRequired = errors.New("player required")
var ErrTeamRequired = errors.New("team required")
