package AppErr

import (
	"errors"
)

var (
	ErrNullField         = errors.New("This Method needs a valid field")
	ErrInvalidArgument   = errors.New("This method needs a valid argument")
	ErrUserNotFound      = errors.New("User not found")
	ErrInformedIncorrect = errors.New("You need to inform the necessary data")
)
