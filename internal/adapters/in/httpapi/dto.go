package httpapi

import (
	"time"

	"github.com/JoelTeoGom/login-service/internal/domain"
)

// DTOs: the JSON shape lives here, not in the domain.
// If the API contract changes, domain never notices.

type registerRequest struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

type updateMeRequest struct {
	Name string `json:"name"`
}

type userResponse struct {
	ID        int64     `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

// toUserResponse never exposes PasswordHash.
func toUserResponse(u *domain.User) userResponse {
	return userResponse{ID: u.ID, Email: u.Email, Name: u.Name, CreatedAt: u.CreatedAt}
}
