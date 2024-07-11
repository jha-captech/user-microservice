package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type Configuration struct {
	Env      string `env:"ENV"`
	Database struct {
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
	cfg := Configuration{}
	if err := env.Parse(&cfg); err != nil {
		return Configuration{}, fmt.Errorf("[in config.New]: %w", err)
	}
	return cfg, nil
}
