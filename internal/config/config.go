package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	DatabaseURL string
	JWTSecret   string
}

func LoadConfig() *Config {
	// Attempt to load .env silently if exists
	_ = godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is required in env")
	}

	jwt := os.Getenv("JWT_SECRET")
	if jwt == "" {
		log.Fatal("JWT_SECRET is required in env")
	}

	return &Config{
		Port:        port,
		DatabaseURL: dbURL,
		JWTSecret:   jwt,
	}
}
