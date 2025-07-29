package verify

import (
	"auth_test/internal/service"
	"bytes"
	"encoding/json"
	"errors"
	"github.com/rs/zerolog/log"
	"net/http"
	"strings"
)

type VerifyHandler struct {
	UserService service.UserService
}

type verifyResponse struct {
	Token string `json:"token"`
}

func (h *VerifyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.Verify(w, r)
}

func (h *VerifyHandler) Verify(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		log.Error().Msg("Invalid authorization header")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	token := strings.TrimPrefix(authHeader, "Bearer ")
	newToken, err := h.UserService.RefreshToken(token)
	if errors.Is(err, service.ErrTokenGeneration) {
		log.Error().Err(service.ErrTokenGeneration).Msg("Token generation failed")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	if err != nil {
		log.Error().Err(err).Msg("Error refreshing token")
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := verifyResponse{newToken}

	var buff bytes.Buffer
	if err := json.NewEncoder(&buff).Encode(response); err != nil {
		log.Error().Err(err).Msg("Error encoding response to JSON")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(buff.Bytes())
}
