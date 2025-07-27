package verify

import (
	"auth_test/internal/service"
	"bytes"
	"encoding/json"
	"net/http"
	"strings"
)

type VerifyHandler struct {
	UserService service.UserService
}

type VerifyResponse struct {
	Token string `json:"token"`
}

func (h *VerifyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.Verify(w, r)
}

func (h *VerifyHandler) Verify(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	token := strings.TrimPrefix(authHeader, "Bearer ")
	newToken, err := h.UserService.RefreshToken(token)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := VerifyResponse{newToken}

	var buff bytes.Buffer
	if err := json.NewEncoder(&buff).Encode(response); err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(buff.Bytes())
}
