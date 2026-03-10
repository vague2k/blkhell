package services

import "errors"

var (
	ErrDb       = errors.New("Database error, contact admin if issue occurs")
	ErrInternal = errors.New("Internal server error")
)
