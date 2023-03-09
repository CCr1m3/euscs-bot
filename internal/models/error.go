package models

import "fmt"

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
	return fmt.Sprintf("User %s is already linked", e.UserID)
}

type UsernameAlreadyLinkedError struct {
	Username string
}

func (e *UsernameAlreadyLinkedError) Error() string {
	return fmt.Sprintf("Username %s is already linked", e.Username)
}

type DBNotFoundError struct {
}

func (e *DBNotFoundError) Error() string {
	return fmt.Sprintf("Not found in db")
}
