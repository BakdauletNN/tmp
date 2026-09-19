package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DataBaseURL string
	Port        string
	JWT_SECRET string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env file not found")
	}

	config := &Config{
		DataBaseURL: os.Getenv("DB_URL"),
		Port:        os.Getenv("PORT"),
		JWT_SECRET:  os.Getenv("JWT_SECRET"),
	}

	return config, nil
}
