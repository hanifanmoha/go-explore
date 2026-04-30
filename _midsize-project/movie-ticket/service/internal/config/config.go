package config

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() error {
	err := godotenv.Load()
	if err != nil {
		slog.Info(".env file not found! Using environment variables instead.")
	}
	return err
}

func GetEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func GetDbUrl() string {
	return GetEnv("DB_URL", "")
}

func GetPort() string {
	return GetEnv("PORT", "6001")
}
