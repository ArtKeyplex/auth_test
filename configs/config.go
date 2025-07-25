package configs

import (
	"flag"
	"os"
)

type Config struct {
	AppName   string
	JwtSecret []byte
	Port      string
}

func LoadConfig() (*Config, error) {
	defaultAppName := "AuthServiceBy_CHEREP"
	defaultJwtSecret := "SuperSecretKey"
	defaultPort := "8080"

	var (
		appName   = flag.String("appName", "", "App Name")
		jwtSecret = flag.String("jwtSecret", "", "Jwt Secret")
		port      = flag.String("port", "", "Port")
	)

	flag.Parse()

	if envAppName := os.Getenv("APP_NAME"); envAppName != "" {
		*appName = envAppName
	}
	if *appName == "" {
		*appName = defaultAppName
	}

	if envJwtSecret := os.Getenv("JWT_SECRET"); envJwtSecret != "" {
		*jwtSecret = envJwtSecret
	}
	if *jwtSecret == "" {
		*jwtSecret = defaultJwtSecret
	}

	if envPort := os.Getenv("PORT"); envPort != "" {
		*port = envPort
	}
	if *port == "" {
		*port = defaultPort
	}

	config := &Config{AppName: *appName, JwtSecret: []byte(*jwtSecret), Port: *port}
	return config, nil
}
