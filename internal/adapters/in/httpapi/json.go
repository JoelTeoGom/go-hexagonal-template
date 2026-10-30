package httpapi

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/JoelTeoGom/login-service/internal/domain"
)

const maxBodyBytes = 1 << 20 // 1 MB

func decodeJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	r.Body = http.MaxBytesReader(w, r.Body, maxBodyBytes)
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil {
		return fmt.Errorf("%w: invalid JSON body", domain.ErrInvalidInput)
	}
	return nil
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
