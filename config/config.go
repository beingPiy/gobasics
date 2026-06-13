package config

import (
	"log"
	"os"
	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
	Port        string
	JWTSecret   string
}

func Load() Config {
	// Load environment variables from .env file
	if err := godotenv.Load() ; err != nil {
		log.Println("No .env file found, using environment variables")
	}
	return Config{
		DatabaseURL: mustGet("DATABASE_URL"),
		Port:        getOrDefault("PORT", "3000"),
		JWTSecret:   mustGet("JWT_SECRET"),
	}
}

func mustGet(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Environment variable %q is required but not set", key)
	}

	return value
}

func getOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}