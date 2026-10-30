package ports

import (
	"context"

	"github.com/JoelTeoGom/login-service/internal/domain"
)

// Outbound ports: what the application NEEDS.
// Implemented by outbound adapters (MySQL, bcrypt, JWT...).

type UserRepository interface {
	// Create stores the user and assigns its ID and CreatedAt.
	// It returns domain.ErrEmailTaken if the email already exists.
	Create(ctx context.Context, u *domain.User) error
	// FindByID and FindByEmail return domain.ErrNotFound if the user does not exist.
	FindByID(ctx context.Context, id int64) (*domain.User, error)
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	Update(ctx context.Context, u *domain.User) error
	// Delete returns domain.ErrNotFound if the user does not exist.
	Delete(ctx context.Context, id int64) error
}

type PasswordHasher interface {
	Hash(password string) (string, error)
	Matches(hash, password string) bool
}

type TokenIssuer interface {
	Issue(userID int64) (string, error)
	Verify(token string) (userID int64, err error)
}
