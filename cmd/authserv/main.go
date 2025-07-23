package main

import (
	"auth_test/configs"
	"auth_test/internal/handlers/login"
	"auth_test/internal/handlers/verify"
	"auth_test/internal/service"
	"auth_test/internal/store"
	"fmt"
	"log"
	"net/http"
)

func main() {
	newStore := &store.InMemoryUserStore{
		Users: map[string]*service.User{
			"Никита": &service.User{Username: "Никита", Password: "password123"},
			"Антон":  &service.User{Username: "Антон", Password: "super123pro"},
			"Адольф": &service.User{Username: "Адольф", Password: "always171wet"},
		},
	}

	config, err := configs.LoadConfig()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Имя приложения:", config.AppName)
	fmt.Println("JWTSecret:", string(config.JwtSecret))
	// USER SERVICE
	userService := service.NewUserService(newStore, config.JwtSecret)
	// LOGIN
	loginHandler := &login.LoginHandler{
		UserService: userService,
	}
	http.Handle("/login", loginHandler)
	// VERIFY
	verifyHandler := &verify.VerifyHandler{
		UserService: userService,
	}
	http.Handle("/verify", verifyHandler)

	// Запуск сервера
	fmt.Println("Сервер запущен на :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
