package user

import "errors"

var (
	// Domain validation errors
	ErrInvalidEmail = errors.New("invalid email address")
	ErrInvalidName  = errors.New("invalid name")
	ErrWeakPassword = errors.New("password must be at least 8 characters")

	// Business logic errors
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUnauthorized      = errors.New("unauthorized")
)
