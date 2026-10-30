package domain

import (
	"net/mail"
	"strings"
	"time"
	"unicode/utf8"
)

// User is the service's central entity. It knows nothing about JSON or SQL.
type User struct {
	ID           int64
	Email        string
	Name         string
	PasswordHash string
	CreatedAt    time.Time
}

// NewUser validates and normalizes the data of a new user.
func NewUser(email, name string) (*User, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	name = strings.TrimSpace(name)

	if err := ValidateEmail(email); err != nil {
		return nil, err
	}
	if err := ValidateName(name); err != nil {
		return nil, err
	}
	return &User{Email: email, Name: name}, nil
}

// Rename changes the name, applying the same rules as on creation.
func (u *User) Rename(name string) error {
	name = strings.TrimSpace(name)
	if err := ValidateName(name); err != nil {
		return err
	}
	u.Name = name
	return nil
}

func ValidateEmail(email string) error {
	addr, err := mail.ParseAddress(email)
	if err != nil || addr.Address != email {
		return invalid("email is not valid")
	}
	return nil
}

func ValidateName(name string) error {
	if n := utf8.RuneCountInString(name); n == 0 || n > 100 {
		return invalid("name must be between 1 and 100 characters")
	}
	return nil
}

// ValidatePassword: bcrypt only uses the first 72 bytes, so the length is capped.
func ValidatePassword(password string) error {
	if len(password) < 8 || len(password) > 72 {
		return invalid("password must be between 8 and 72 bytes")
	}
	return nil
}
