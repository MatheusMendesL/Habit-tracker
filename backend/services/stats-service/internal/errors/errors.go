package AppErr

import (
	"errors"
)

var (
	ErrNullField         = errors.New("This Method needs a valid field")
	ErrInvalidArgument   = errors.New("This method needs a valid argument")
	ErrUserNotFound      = errors.New("User not found")
	ErrInternalError     = errors.New("Ocurred an internal error")
	ErrUserStatsNotFound = errors.New("User Stats not found")
)
