package domain

import "errors"

// Domain errors
var (
	ErrInvalidInput          = errors.New("invalid input")
	ErrInvalidCredentials    = errors.New("invalid credentials")
	ErrUserNotFound          = errors.New("user not found")
	ErrEmailExists           = errors.New("email already exists")
	ErrEmailAlreadyExists    = errors.New("email already exists")
	ErrEmailNotFound         = errors.New("email not found")
	ErrUsernameAlreadyExists = errors.New("username already exists")
	ErrUnauthorized          = errors.New("unauthorized")
)