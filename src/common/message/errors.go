package message

import "errors"

var (
	ErrInvalidToken    = errors.New("invalid token")
	ErrInactiveUser    = errors.New("user not activated")
	ErrInvalidInput    = errors.New("invalid input")
	ErrNoActiveSession = errors.New("user has no active session")
	ErrUserNotFound    = errors.New("user not found")
)
