package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DB_URL string
	PORT   string
	// JWT_SECRET string
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		return value
	}
	return fallback
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("couldn't load env file: ", err)
	}

	return &Config{
		DB_URL: getEnv("DB_URL", "postgres://localhost:5432/default_db"),
		PORT:   getEnv("PORT", "8080"),
	}
}
