package httpapi

import (
	"errors"
	"net/http"

	"github.com/JoelTeoGom/login-service/internal/domain"
)

// writeError is the ONLY place where a domain error becomes HTTP.
func (s *Server) writeError(w http.ResponseWriter, r *http.Request, err error) {
	status, msg := http.StatusInternalServerError, "internal error"

	switch {
	case errors.Is(err, domain.ErrInvalidInput):
		status, msg = http.StatusBadRequest, err.Error()
	case errors.Is(err, domain.ErrEmailTaken):
		status, msg = http.StatusConflict, err.Error()
	case errors.Is(err, domain.ErrInvalidCredentials),
		errors.Is(err, domain.ErrUnauthorized):
		status, msg = http.StatusUnauthorized, err.Error()
	case errors.Is(err, domain.ErrNotFound):
		status, msg = http.StatusNotFound, err.Error()
	default:
		// 500s are logged in detail, but nothing leaks to the client.
		s.log.Error("request failed", "method", r.Method, "path", r.URL.Path, "err", err)
	}

	writeJSON(w, status, map[string]string{"error": msg})
}
