package config

import (
	"fmt"
	"os"

	errors "github.com/Faner201/sc_links/internal/error"
	"github.com/joho/godotenv"
)

type Config struct {
	BaseURL                 string
	Host                    string
	Port                    string
	TelegramContactUsername string
	Github                  GithubConfig
	DB                      DBConfig
	Auth                    AuthConfig
}

type DBConfig struct {
	URI      string
	Database string
}

type AuthConfig struct {
	JWTSecretKey     string
	AllowedGithubOrg string
}

type GithubConfig struct {
	ClientID     string
	ClientSecret string
}

func (c Config) ListenAddr() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}

func Init() error {
	err := godotenv.Load("../.env")

	if err != nil {
		return errors.ErrReadConfig
	}
	return nil
}

func Get() *Config {
	uriDB := os.Getenv("mongodb_uri")
	nameDB := os.Getenv("mongodb_name")
	baseURL := os.Getenv("base_url")
	host := os.Getenv("host")
	port := os.Getenv("port")
	jwt := os.Getenv("jwt_secret_key")
	ghOrg := os.Getenv("allowed_github_org")
	telegramCont := os.Getenv("telegram_contact_username")
	clientID := os.Getenv("github_client_id")
	clientSecret := os.Getenv("github_client_secret")

	return &Config{
		BaseURL:                 baseURL,
		Host:                    host,
		Port:                    port,
		TelegramContactUsername: telegramCont,
		DB: DBConfig{
			URI:      uriDB,
			Database: nameDB,
		},
		Auth: AuthConfig{
			JWTSecretKey:     jwt,
			AllowedGithubOrg: ghOrg,
		},
		Github: GithubConfig{
			ClientID:     clientID,
			ClientSecret: clientSecret,
		},
	}
}
