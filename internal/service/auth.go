package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/JoelTeoGom/login-service/internal/domain"
	"github.com/JoelTeoGom/login-service/internal/ports"
)

// AuthService implements ports.AuthService. It only knows interfaces:
// it has no idea whether MySQL, bcrypt or JWT sit behind them.
type AuthService struct {
	users  ports.UserRepository
	hasher ports.PasswordHasher
	tokens ports.TokenIssuer
}

var _ ports.AuthService = (*AuthService)(nil)

func NewAuthService(users ports.UserRepository, hasher ports.PasswordHasher, tokens ports.TokenIssuer) *AuthService {
	return &AuthService{users: users, hasher: hasher, tokens: tokens}
}

func (s *AuthService) Register(ctx context.Context, in ports.RegisterInput) (*domain.User, error) {
	u, err := domain.NewUser(in.Email, in.Name)
	if err != nil {
		return nil, err
	}
	if err := domain.ValidatePassword(in.Password); err != nil {
		return nil, err
	}

	hash, err := s.hasher.Hash(in.Password)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}
	u.PasswordHash = hash

	if err := s.users.Create(ctx, u); err != nil {
		return nil, err
	}
	return u, nil
}

func (s *AuthService) Login(ctx context.Context, in ports.LoginInput) (string, error) {
	email := strings.ToLower(strings.TrimSpace(in.Email))

	u, err := s.users.FindByEmail(ctx, email)
	if errors.Is(err, domain.ErrNotFound) {
		// Same error as a wrong password: never reveal which emails exist.
		return "", domain.ErrInvalidCredentials
	}
	if err != nil {
		return "", err
	}

	if !s.hasher.Matches(u.PasswordHash, in.Password) {
		return "", domain.ErrInvalidCredentials
	}

	token, err := s.tokens.Issue(u.ID)
	if err != nil {
		return "", fmt.Errorf("issue token: %w", err)
	}
	return token, nil
}

func (s *AuthService) VerifyToken(_ context.Context, token string) (int64, error) {
	id, err := s.tokens.Verify(token)
	if err != nil {
		return 0, domain.ErrUnauthorized
	}
	return id, nil
}
