package ports

import (
	"context"

	"github.com/JoelTeoGom/login-service/internal/domain"
)

// Inbound ports: what the application OFFERS.
// Used by inbound adapters (HTTP, gRPC, CLI...).

type RegisterInput struct {
	Email    string
	Name     string
	Password string
}

type LoginInput struct {
	Email    string
	Password string
}

type AuthService interface {
	Register(ctx context.Context, in RegisterInput) (*domain.User, error)
	Login(ctx context.Context, in LoginInput) (token string, err error)
	VerifyToken(ctx context.Context, token string) (userID int64, err error)
}

type UserService interface {
	Get(ctx context.Context, id int64) (*domain.User, error)
	UpdateName(ctx context.Context, id int64, name string) (*domain.User, error)
	Delete(ctx context.Context, id int64) error
}
