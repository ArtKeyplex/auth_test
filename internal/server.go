package internal

import (
	"auth_test/internal/handlers/login"
	"auth_test/internal/service"
	"net/http"
	"time"
)

func NewServer(loginHandler http.HandlerFunc, verifyHandler http.HandlerFunc) *http.Server {
	mux := http.NewServeMux()

	// Регистрируем маршруты
	mux.HandleFunc("/login", loginHandler)
	mux.HandleFunc("/verify", verifyHandler)

	srv := &http.Server{
		Addr:         ":8080", // порт, на котором слушаем
		Handler:      mux,     // наш mux с маршрутами
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	return srv
}

func NewLoginHandler(userService *service.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	UserService:
		userService
	}
}

func NewVerifyHandler(userService *service.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Ваша логика verify, использующая userService
	}
}
