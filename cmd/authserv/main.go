package main

import (
	"auth_test/configs"
	"auth_test/internal"
	"auth_test/internal/handlers/login"
	"auth_test/internal/handlers/verify"
	"auth_test/internal/service"
	"auth_test/internal/store"
	"fmt"
	"log"
)

func main() {
	newStore := store.NewInMemoryUserStore()

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
	// VERIFY
	verifyHandler := &verify.VerifyHandler{
		UserService: userService,
	}
	srv := internal.NewServer(loginHandler.ServeHTTP, verifyHandler.ServeHTTP)
	fmt.Println("Сервер слушает на", srv.Addr)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}

}
