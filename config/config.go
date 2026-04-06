package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL        string
	ServerPort         string
	FRONTEND_URL       string
	ACCESS_SECRET_KEY  string
	REFRESH_SECRET_KEY string
	RequestTimeout     time.Duration
}

func Load() (Config, error) {
	fmt.Println("Loading global configuration...")
	err := godotenv.Load()
	if err != nil {
		return Config{}, err
	}

	config := Config{
		DatabaseURL:        getEnvOrDefault("DATABASE_URL", "postgres://user:password@localhost:5432/talky_space"),
		ServerPort:         getEnvOrDefault("SERVER_PORT", "8080"),
		FRONTEND_URL:       getEnvOrDefault("FRONTEND_URL", "http://localhost:3000"),
		ACCESS_SECRET_KEY:  getEnvOrDefault("ACCESS_TOKEN_SECRET", ""),
		REFRESH_SECRET_KEY: getEnvOrDefault("REFRESH_TOKEN_SECRET", ""),
		RequestTimeout:     time.Duration(getEnvIntOrDefault("REQUEST_TIMEOUT_SECONDS", 5)) * time.Second,
	}

	return config, nil
}

func getEnvOrDefault(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func getEnvIntOrDefault(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
}
