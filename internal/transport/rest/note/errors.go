package note

import "errors"

var (
	ErrInvalidPageParam  = errors.New("invalid page param")
	ErrInvalidLimitParam = errors.New("invalid limit param")
)
