package config


import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct{
	DataBaseURL string
	Port string
}

func LoadConfig() (*Config, error){
	var err error = godotenv.LoadConfig()
	if err != nil{
		log.Println(".env file not found")
	}
	var config *Config = &Config{
		DataBaseURL: os.Getenv("DB_URL"),
		Port: os.Getenv("Port"),
	}
}