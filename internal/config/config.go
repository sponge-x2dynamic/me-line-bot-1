package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port              string
	LineChannelSecret string
	LineAccessToken   string
	DBHost            string
	DBPort            string
	DBUser            string
	DBPassword        string
	DBName            string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, reading from environment")
	}

	return &Config{
		Port:              getEnv("PORT", "8080"),
		LineChannelSecret: getEnv("LINE_CHANNEL_SECRET", ""),
		LineAccessToken:   getEnv("LINE_CHANNEL_ACCESS_TOKEN", ""),
		DBHost:            getEnv("DB_HOST", "localhost"),
		DBPort:            getEnv("DB_PORT", "3306"),
		DBUser:            getEnv("DB_USER", "mebot"),
		DBPassword:        getEnv("DB_PASSWORD", ""),
		DBName:            getEnv("DB_NAME", "mebot_db"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
