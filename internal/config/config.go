package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DataBaseURL  string
	Port         string
	JWT_SECRET   string
	RedisURL     string
	SMTPHost     string
	SMTPPort     string
	SMTPUser     string
	SMTPPassword string
	SMTPFrom     string
	WebBaseURL   string
}

func LoadConfig() (*Config, error) {
	_ = godotenv.Load()

	config := &Config{
		DataBaseURL:  os.Getenv("DB_URL"),
		Port:         os.Getenv("PORT"),
		JWT_SECRET:   os.Getenv("JWT_SECRET"),
		RedisURL:     os.Getenv("REDIS_URL"),
		SMTPHost:     os.Getenv("SMTP_HOST"),
		SMTPPort:     envOrDefault("SMTP_PORT", "587"),
		SMTPUser:     os.Getenv("SMTP_USER"),
		SMTPPassword: os.Getenv("SMTP_PASSWORD"),
		SMTPFrom:     os.Getenv("SMTP_FROM"),
		WebBaseURL:   envOrDefault("WEB_BASE_URL", "http://localhost:5173"),
	}

	return config, nil
}

func envOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
