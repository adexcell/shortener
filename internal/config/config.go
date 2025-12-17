package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseDSN string
	HTTPServer  string
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on environment variables")
	}

	return &Config{
		DatabaseDSN: getEnv("DATABASE_DSN", "postgres://admin:password@localhost:5432/shortener?sslmode=disable"),
		HTTPServer: getEnv("HTTP_SERVER_ADDRESS", ":8080"),
	}
}

func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}