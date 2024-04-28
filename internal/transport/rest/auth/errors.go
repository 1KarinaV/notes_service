package auth

import "errors"

var (
	ErrInvalidCast     = errors.New("invalid type cast")
	ErrInvalidPassword = errors.New("invalid password")
)
