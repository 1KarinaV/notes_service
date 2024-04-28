package services

import "errors"

var (
	ErrNotFound           = errors.New("not found")
	ErrEmptyClaims        = errors.New("empty claims")
	ErrNoteUpdateExceeded = errors.New("update time exceeded")
)
