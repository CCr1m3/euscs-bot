package models

import (
	"errors"
	"fmt"
)

type DBError struct {
	Err error
}

func (e *DBError) Error() string {
	return fmt.Sprintf("database error: %s", e.Err.Error())
}

func (e *DBError) Unwrap() error {
	return e.Err
}

type RankUpdateUsernameError struct {
	Username string
}

func (e *RankUpdateUsernameError) Error() string {
	return fmt.Sprintf("rank update username %s invalid", e.Username)
}

type RankUpdateTooFastError struct {
	UserID string
}

func (e *RankUpdateTooFastError) Error() string {
	return fmt.Sprintf("rank update too fast for user %s", e.UserID)
}

type NotLinkedError struct {
	UserID string
}

func (e *NotLinkedError) Error() string {
	return fmt.Sprintf("player %s has not linked an omega strikers account", e.UserID)
}

type UserAlreadyLinkedError struct {
	UserID string
}

func (e *UserAlreadyLinkedError) Error() string {
	return fmt.Sprintf("user %s is already linked", e.UserID)
}

type UsernameAlreadyLinkedError struct {
	Username string
}

func (e *UsernameAlreadyLinkedError) Error() string {
	return fmt.Sprintf("username %s is already linked", e.Username)
}

type DBNotFoundError struct {
}

func (e *DBNotFoundError) Error() string {
	return fmt.Sprintf("not found in db")
}

type TeamIsFullError struct {
}

func (e *TeamIsFullError) Error() string {
	return fmt.Sprintf("team is full")
}

type UserAlreadyInTeam struct {
}

func (e *UserAlreadyInTeam) Error() string {
	return fmt.Sprintf("user is already in a team")
}

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
