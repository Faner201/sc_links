package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	BaseURL string
	Host    string
	Port    string
	DB      DBConfig
}

type DBConfig struct {
	URI      string
	Database string
}

func (c Config) ListenAddr() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}

func Get() *Config {
	err := godotenv.Load("../.env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	uriDB := os.Getenv("MONGODB_URI")
	nameDB := os.Getenv("MONGODB_NAME")
	baseURL := os.Getenv("BASE_URL")
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")

	return &Config{
		BaseURL: baseURL,
		Host:    host,
		Port:    port,
		DB: DBConfig{
			URI:      uriDB,
			Database: nameDB,
		},
	}
}
