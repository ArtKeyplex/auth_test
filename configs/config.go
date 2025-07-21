package configs

import (
	"flag"
	"os"
)

type Config struct {
	AppName string
}

func LoadConfig() (*Config, error) {
	defaultAppName := "AuthServiceBy_CHEREP"

	var (
		appName = flag.String("appName", "", "Имя приложения")
	)

	flag.Parse()

	if envAppName := os.Getenv("APP_NAME"); envAppName != "" {
		*appName = envAppName
	}

	if *appName == "" {
		*appName = defaultAppName
	}

	config := &Config{AppName: *appName}
	return config, nil
}
