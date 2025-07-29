package configs

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog/log"
)

type Config struct {
	AppName   string `envconfig:"APP_NAME" default:"AuthService"`
	JwtSecret string `envconfig:"JWT_SECRET"`
	Port      int    `envconfig:"PORT"`
	DSN       string `envconfig:"DSN"`
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Error().Err(err).Msg("Error loading .env file")
	}
	var cfg Config

	// Считываем данные из .env файла в переменные
	if err := envconfig.Process("", &cfg); err != nil {
		log.Error().Err(err).Msg("Error reading .env file into config")
		return nil, err
	}

	return &cfg, nil
}
