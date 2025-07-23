package login

import (
	"auth_test/internal/service"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"
)

type LoginHandler struct {
	UserService service.UserService
}

func (h *LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Basic ") {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	payload := strings.TrimPrefix(authHeader, "Basic ")
	decoded, err := base64.StdEncoding.DecodeString(payload)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	credentials := strings.SplitN(string(decoded), ":", 2)
	if len(credentials) != 2 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	username, password := credentials[0], credentials[1]

	valid, err := h.UserService.ValidateCredentials(username, password)
	if err != nil || !valid {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	token, err := h.UserService.GenerateToken(username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Если дошли сюда значит все наши кейсы прошли, ошибок нет,
	// токен сгенерирован, введенные данные валидны, тогда
	// отправляем ответ в json и заголовке Auth bearer+token

	// Устанавливаем заголовок Authorization с Bearer токеном
	w.Header().Set("Authorization", "Bearer "+token)
	// Возвращаем статус 200
	w.WriteHeader(http.StatusOK)

	// Возвращаем json
	response := map[string]string{"token": token}
	json.NewEncoder(w).Encode(response)
	return
}
