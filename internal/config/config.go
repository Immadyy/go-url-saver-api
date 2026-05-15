package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBURL string
	PORT  string
}

// Excellent helper function structure
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		return value
	}
	return fallback
}

func LoadConfig() *Config {
	// Suppress filesystem errors safely for container layers
	_ = godotenv.Load()

	// Utilize your helper function for clean fallbacks
	port := getEnv("PORT", "8080")
	dbURL := getEnv("DB_URL", "")

	if dbURL == "" {
		log.Println("Warning: DB_URL environment variable is empty")
	}

	return &Config{
		DBURL: dbURL,
		PORT:  port,
	}
}
