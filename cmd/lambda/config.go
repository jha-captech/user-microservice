package main

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/joho/godotenv"

	"github.com/spf13/viper"
)

type configuration struct {
	Env      string
	Database struct {
		Name            string
		User            string
		Password        string
		Host            string
		Port            string
		ConnectionRetry int
	}
	HTTP struct {
		Domain string
		Port   string
	}
}

func mustNewConfiguration() configuration {
	loadEnv()

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath("./")

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Sprintf("unable to decode config file, %v", err))
	}

	config := configuration{}
	if err := viper.Unmarshal(&config); err != nil {
		panic(fmt.Sprintf("unable to decode config file into configuration, %v", err))
	}

	return config
}

func loadEnv() {
	if err := godotenv.Load(); err != nil {
		slog.Info("DOTENV: No .env file found or error reading file.")
		return
	}
	slog.Info("DOTENV: .env found.")
}
