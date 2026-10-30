package httpapi

import (
	"log/slog"
	"net/http"

	"github.com/JoelTeoGom/login-service/internal/ports"
)

// Server is the inbound HTTP adapter. It only depends on the inbound
// ports: it does not know MySQL, bcrypt or JWT exist.
type Server struct {
	auth  ports.AuthService
	users ports.UserService
	log   *slog.Logger
}

// NewHandler mounts the routes and returns an http.Handler ready for the http.Server.
func NewHandler(auth ports.AuthService, users ports.UserService, log *slog.Logger) http.Handler {
	s := &Server{auth: auth, users: users, log: log}

	mux := http.NewServeMux()

	// Public
	mux.Handle("POST /auth/register", s.handle(s.register))
	mux.Handle("POST /auth/login", s.handle(s.login))
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Protected with JWT
	mux.Handle("GET /me", s.requireAuth(s.handle(s.getMe)))
	mux.Handle("PATCH /me", s.requireAuth(s.handle(s.updateMe)))
	mux.Handle("DELETE /me", s.requireAuth(s.handle(s.deleteMe)))

	return mux
}

// handlerFunc is like an http.HandlerFunc but returns an error.
// That way each handler does "return err" and status mapping lives in one place.
type handlerFunc func(w http.ResponseWriter, r *http.Request) error

func (s *Server) handle(h handlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			s.writeError(w, r, err)
		}
	})
}
