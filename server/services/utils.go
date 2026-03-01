package services

import "errors"

var (
	ErrDb               = errors.New("Database: Error, contact admin if issue occurs")
	ErrNoSession        = errors.New("Session: doesn't exist")
	ErrFileNotSupported = errors.New("File: format is not supported")
	ErrNoPath           = errors.New("Path: doesn't exist")
	ErrInternal         = errors.New("Internal: server error")
)
