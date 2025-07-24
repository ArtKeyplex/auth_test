package verify

import (
	"auth_test/internal/service"
	"encoding/json"
	"net/http"
	"strings"
)

type VerifyHandler struct {
	UserService service.UserService
}

func (h *VerifyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("заголовок пустой или неверный"))
		return
	}
	token := strings.TrimPrefix(authHeader, "Bearer ")
	newToken, err := h.UserService.RefreshToken(token)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("ошибка в методе refresh token\n" + err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	response := map[string]string{"refreshed token": newToken}
	json.NewEncoder(w).Encode(response)
	return
}
