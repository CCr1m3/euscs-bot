package static

import (
	"errors"
	"fmt"
)

func ErrDB(err error) error {
	return fmt.Errorf("database error: %w", err)
}

var ErrUsernameInvalid = errors.New("username invalid")
var ErrUserNotLinked = errors.New("user not linked")
var ErrUserAlreadyLinked = errors.New("user is already linked")
var ErrUsernameAlreadyLinked = errors.New("username is already linked")
var ErrNotFound = errors.New("not found")
var ErrTeamFull = errors.New("team full")
var ErrUserAlreadyInTeam = errors.New("user already in team")
var ErrTeamnameTaken = errors.New("team name is taken")
var ErrNotTeamOwner = errors.New("not team owner")
var ErrOwnerNotInTeam = errors.New("owner not in team")
var ErrCorestrikeNotFound = errors.New("username not found in corestrike")
var ErrRankUpdateTooFast = errors.New("rank update too fast")
var ErrDiscordIDRequired = errors.New("discord id required")
