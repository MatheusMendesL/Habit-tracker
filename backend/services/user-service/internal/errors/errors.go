package errors

import (
	"errors"
)

var (
	ErrNullField       = errors.New("This Method needs a valid field")
	ErrInvalidArgument = errors.New("This method needs a valid argument")
	ErruUserNotFound   = errors.New("User not found")
)
