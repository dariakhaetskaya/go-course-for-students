package app

import (
	"errors"
)

var (
	ErrNotAuthor = errors.New("only author can edit the ad")
	ErrNotFound  = errors.New("requested ad not found")
	ErrInvalid   = errors.New("ad content is invalid")
)
