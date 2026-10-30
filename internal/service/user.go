package service

import (
	"context"

	"github.com/JoelTeoGom/login-service/internal/domain"
	"github.com/JoelTeoGom/login-service/internal/ports"
)

// UserService implements ports.UserService: CRUD on the user's own account.
type UserService struct {
	users ports.UserRepository
}

var _ ports.UserService = (*UserService)(nil)

func NewUserService(users ports.UserRepository) *UserService {
	return &UserService{users: users}
}

func (s *UserService) Get(ctx context.Context, id int64) (*domain.User, error) {
	return s.users.FindByID(ctx, id)
}

func (s *UserService) UpdateName(ctx context.Context, id int64, name string) (*domain.User, error) {
	u, err := s.users.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := u.Rename(name); err != nil {
		return nil, err
	}
	if err := s.users.Update(ctx, u); err != nil {
		return nil, err
	}
	return u, nil
}

func (s *UserService) Delete(ctx context.Context, id int64) error {
	return s.users.Delete(ctx, id)
}
