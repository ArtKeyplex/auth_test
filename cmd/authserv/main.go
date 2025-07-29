package main

import (
	"auth_test/configs"
	"auth_test/internal"
	"auth_test/internal/service"
	"auth_test/internal/store"
	"context"
	"github.com/rs/zerolog"
	"os"
)

func main() {
	ctx := context.Background()
	log := zerolog.New(os.Stdout).With().Timestamp().Logger()

	// Загружаем конфиг
	config, err := configs.LoadConfig()
	if err != nil {
		log.Error().Err(err).Msg("failed to load config")
		return
	}
	log.Info().Msg("config loaded")

	// Подключаемся к PSQL
	db, err := store.InitDb(ctx, config)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database")
	}
	log.Info().Msg("connected to database")

	// Инициализируем наше хранилище
	newStore := store.NewPostgresUserStore(db)
	// Создаем тип, удовлетворяющий интерфейсу UserService
	userService := service.NewUserService(newStore, []byte(config.JwtSecret))

	// Создаем обработчики http-запросов
	loginHandler := internal.NewHandler("login", userService)
	verifyHandler := internal.NewHandler("verify", userService)

	// Создаем сервер
	srv := internal.NewServer(loginHandler.ServeHTTP, verifyHandler.ServeHTTP)
	log.Info().Str("PORT", srv.Addr).Msg("starting server")
	// newStore.AddUser(ctx, db, "Никита", "password123")
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal().Err(err).Msg("server stopped")
	}

}
