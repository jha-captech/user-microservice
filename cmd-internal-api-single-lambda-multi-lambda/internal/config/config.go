package config

import (
	"fmt"
	"log/slog"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Configuration struct {
	Env        string     `env:"ENV,required"`
	LogLevel   slog.Level `env:"LOG_LEVEL,required"`
	UseSwagger bool       `env:"USE_SWAGGER" envDefault:"false"`
	Database   struct {
		Name            string `env:"DATABASE_NAME"`
		User            string `env:"DATABASE_USER"`
		Password        string `env:"DATABASE_PASSWORD"`
		Host            string `env:"DATABASE_HOST"`
		Port            string `env:"DATABASE_PORT"`
		ConnectionRetry int    `env:"DATABASE_CONNECTION_RETRY"`
	}
	HTTP struct {
		Domain              string `env:"HTTP_DOMAIN"`
		Port                string `env:"HTTP_PORT"`
		ShutdownGracePeriod int    `env:"HTTP_SHUTDOWN_GRACE_PERIOD"`
	}
}

func New() (Configuration, error) {
	_ = godotenv.Load()

	cfg, err := env.ParseAs[Configuration]()
	if err != nil {
		return Configuration{}, fmt.Errorf("[in config.New]: %w", err)
	}

	return cfg, nil
}
