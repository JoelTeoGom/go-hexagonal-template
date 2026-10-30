package httpapi

import (
	"context"
	"net/http"
	"strings"

	"github.com/JoelTeoGom/login-service/internal/domain"
)

type ctxKey struct{}

// requireAuth validates "Authorization: Bearer <jwt>" and puts the userID in the context.
func (s *Server) requireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw, ok := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer ")
		if !ok || raw == "" {
			s.writeError(w, r, domain.ErrUnauthorized)
			return
		}

		userID, err := s.auth.VerifyToken(r.Context(), raw)
		if err != nil {
			s.writeError(w, r, err)
			return
		}

		ctx := context.WithValue(r.Context(), ctxKey{}, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func userIDFrom(ctx context.Context) int64 {
	id, _ := ctx.Value(ctxKey{}).(int64)
	return id
}
