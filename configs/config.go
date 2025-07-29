package configs

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	AppName   string `envconfig:"APP_NAME" default:"AuthService"`
	JwtSecret string `envconfig:"JWT_SECRET" default:"SecretKey"`
	Port      int    `envconfig:"PORT" default:"8081"`
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Нет файла .env или ошибка чтения")
	}
	var cfg Config

	if err := envconfig.Process("", &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
