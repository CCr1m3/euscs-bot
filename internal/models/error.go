package models

import "fmt"

type DBError struct {
	Err error
}

func (e *DBError) Error() string {
	return fmt.Sprintf("database error: %s" + e.Err.Error())
}

func (e *DBError) Unwrap() error {
	return e.Err
}

type RankUpdateUsernameError struct {
	Username string
}

func (e *RankUpdateUsernameError) Error() string {
	return fmt.Sprintf("rank update username %s invalid: ", e.Username)
}
