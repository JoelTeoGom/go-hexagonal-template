package domain

import (
	"errors"
	"fmt"
)

// Domain errors. Adapters translate them: the HTTP one into status codes,
// the MySQL one turns its own errors into these.
var (
	ErrNotFound           = errors.New("not found")
	ErrEmailTaken         = errors.New("email already registered")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrInvalidInput       = errors.New("invalid input")
)

func invalid(msg string) error {
	return fmt.Errorf("%w: %s", ErrInvalidInput, msg)
}
