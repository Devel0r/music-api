package config

import (
	"errors"
	"log"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

var ErrFailedParseEnv = errors.New("error, failed parse env")

type Config struct {
	DB struct {
		Host     string `envconfig:"DB_HOST" default:"localhost"`
		Port     string `envconfig:"DB_PORT" default:"5432"`
		User     string `envconfig:"DB_USER" default:"postgres"`
		Password string `envconfig:"DB_PASSWORD" default:"postgres"`
		Name     string `envconfig:"DB_NAME" default:"music_api"`
	}
	Server struct {
		Port string `envconfig:"SERVER_PORT" default:"8080"`
	}
	ExternalAPI string `envconfig:"EXTERNAL_API_URL"`
}

func InitConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, ErrFailedParseEnv
	}

	return &cfg, nil
}
