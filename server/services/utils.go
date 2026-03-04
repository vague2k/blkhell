package services

import "errors"

var (
	ErrDb               = errors.New("Database error, contact admin if issue occurs")
	ErrNoSession        = errors.New("Session doesn't exist")
	ErrFileNotSupported = errors.New("The file format is not supported")
	ErrFileTooBig       = errors.New("The uploaded file is too big. Please choose an file that's less than 100MB in size")
	ErrNoPath           = errors.New("The path doesn't exist")
	ErrInternal         = errors.New("Internal server error")
)
