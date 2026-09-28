package services

import "errors"

var (
	ErrInvalidInput       = errors.New("invalid input")
	ErrInvalidCategory    = errors.New("invalid category")
	ErrInvalidSearchQuery = errors.New("invalid search query")
	ErrNotFound           = errors.New("not found")
	ErrEmailTaken         = errors.New("email already registered")
	ErrInvalidCredentials = errors.New("invalid email or password")
)
