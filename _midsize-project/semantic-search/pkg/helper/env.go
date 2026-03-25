package helper

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}

func GetEnv(key string) string {
	return os.Getenv(key)
}

func GetGenAIAPIKey() string {
	return GetEnv("GENAI_API_KEY")
}

func GetDBURL() string {
	return GetEnv("DB_URL")
}
