package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

type Config struct {
	Env      string         `env:"ENV" env-default:"local"`
	Database DatabaseConfig `env:"DATABASE" env-required:"true"`
	Server   ServerConfig   `env:"HTTP_SERVER" env-required:"true"`
}

type DatabaseConfig struct {
	Host     string `env:"DATABASE_HOST"`
	Port     int    `env:"DATABASE_PORT"`
	Name     string `env:"DATABASE_NAME"`
	User     string `env:"DATABASE_USER"`
	Password string `env:"DATABASE_PASSWORD"`
}

type ServerConfig struct {
	Address     string        `env:"HTTP_SERVER_ADDRESS" env-default:"localhost:8080"`
	Timeout     time.Duration `env:"HTTP_SERVER_TIMEOUT" env-default:"4s"`
	IdleTimeout time.Duration `env:"HTTP_SERVER_IDLE_TIMEOUT" env-default:"60s"`
}

func MustLoad() *Config {
	configPath := "config/config.env"

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatal("config file doesn't exists: ", err)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatal("failed to read config file: ", err)
	}
	return &cfg
}
