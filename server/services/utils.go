package services

import "errors"

var (
	ErrDb        = errors.New("Database error, try again. Contact admin if issue occurs")
	ErrNoSession = errors.New("Session does not exist")
)
