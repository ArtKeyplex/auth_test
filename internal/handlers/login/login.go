package login

import (
	"auth_test/internal/service"
	"bytes"
	"encoding/json"
	"github.com/rs/zerolog/log"
	"net/http"
)

type LoginHandler struct {
	UserService service.UserService
}

type loginResponse struct {
	Token string `json:"token"`
}

func (h *LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.Login(w, r)
}

func (h *LoginHandler) Login(w http.ResponseWriter, r *http.Request) {
	username, password, ok := r.BasicAuth()
	if !ok {
		log.Error().Msg("Basic auth failed")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	valid, err := h.UserService.ValidateCredentials(username, password)
	// TODO: connection error
	if err != nil {
		log.Error().Err(err).Msg("Error validating username")
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	if !valid {
		log.Error().Msg("Error validating password")
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	token, err := h.UserService.GenerateToken(username)
	if err != nil {
		log.Error().Err(err).Msg("Error generating token")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	// Если дошли сюда значит все наши кейсы прошли, ошибок нет,
	// токен сгенерирован, введенные данные валидны, тогда
	// отправляем ответ в json и заголовке Auth bearer+token

	// Устанавливаем заголовок Authorization с Bearer токеном
	w.Header().Set("Authorization", "Bearer "+token)
	w.Header().Set("Content-Type", "application/json")
	response := loginResponse{Token: token}

	var buff bytes.Buffer
	if err := json.NewEncoder(&buff).Encode(response); err != nil {
		log.Error().Err(err).Msg("Error encoding response to JSON")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(buff.Bytes())
}
