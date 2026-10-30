package httpapi

import (
	"net/http"

	"github.com/JoelTeoGom/login-service/internal/ports"
)

// Every handler does the same thing, in this order:
//   1. decode the request (DTO)
//   2. call the inbound port
//   3. map the response (DTO)
// No business logic here. Errors are returned and writeError translates them.

func (s *Server) register(w http.ResponseWriter, r *http.Request) error {
	var req registerRequest
	if err := decodeJSON(w, r, &req); err != nil {
		return err
	}

	u, err := s.auth.Register(r.Context(), ports.RegisterInput{
		Email:    req.Email,
		Name:     req.Name,
		Password: req.Password,
	})
	if err != nil {
		return err
	}

	writeJSON(w, http.StatusCreated, toUserResponse(u))
	return nil
}

func (s *Server) login(w http.ResponseWriter, r *http.Request) error {
	var req loginRequest
	if err := decodeJSON(w, r, &req); err != nil {
		return err
	}

	token, err := s.auth.Login(r.Context(), ports.LoginInput{Email: req.Email, Password: req.Password})
	if err != nil {
		return err
	}

	writeJSON(w, http.StatusOK, loginResponse{AccessToken: token, TokenType: "Bearer"})
	return nil
}

func (s *Server) getMe(w http.ResponseWriter, r *http.Request) error {
	u, err := s.users.Get(r.Context(), userIDFrom(r.Context()))
	if err != nil {
		return err
	}

	writeJSON(w, http.StatusOK, toUserResponse(u))
	return nil
}

func (s *Server) updateMe(w http.ResponseWriter, r *http.Request) error {
	var req updateMeRequest
	if err := decodeJSON(w, r, &req); err != nil {
		return err
	}

	u, err := s.users.UpdateName(r.Context(), userIDFrom(r.Context()), req.Name)
	if err != nil {
		return err
	}

	writeJSON(w, http.StatusOK, toUserResponse(u))
	return nil
}

func (s *Server) deleteMe(w http.ResponseWriter, r *http.Request) error {
	if err := s.users.Delete(r.Context(), userIDFrom(r.Context())); err != nil {
		return err
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}
