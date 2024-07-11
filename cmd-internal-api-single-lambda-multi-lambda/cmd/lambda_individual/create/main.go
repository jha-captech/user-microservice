package main

import (
	"fmt"
	"log"
	"log/slog"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
	"github.com/caarlos0/env/v11"
	"github.com/go-chi/chi/v5"
	"github.com/jha-captech/user-microservice/internal/database"
	"github.com/jha-captech/user-microservice/internal/handlers"
	"github.com/jha-captech/user-microservice/internal/user"
)

type configuration struct {
	Env      string `env:"ENV"`
	Database struct {
		Name            string `env:"DATABASE_NAME,required"`
		User            string `env:"DATABASE_USER,required"`
		Password        string `env:"DATABASE_PASSWORD,required"`
		Host            string `env:"DATABASE_HOST,required"`
		Port            string `env:"DATABASE_PORT,required"`
		ConnectionRetry int    `env:"DATABASE_CONNECTION_RETRY,required"`
	}
}

func main() {
	if err := run(); err != nil {
		log.Fatalf("Startup failed. err: %v", err)
	}
}

func run() error {
	cfg := configuration{}
	if err := env.Parse(&cfg); err != nil {
		return fmt.Errorf("[in run]: %w", err)
	}

	logger := slog.Default()

	db, err := database.NewDatabase(
		fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
			cfg.Database.Host,
			cfg.Database.User,
			cfg.Database.Password,
			cfg.Database.Name,
			cfg.Database.Port,
		),
		logger,
		cfg.Database.ConnectionRetry,
	)
	if err != nil {
		return fmt.Errorf("[in run]: %w", err)
	}

	defer db.Session.Close()

	r := chi.NewRouter()

	us := user.NewService(db)
	h := handlers.New(logger, us)

	r.Post("/api/user", h.HandleCreateUser())

	lambda.Start(httpadapter.New(r).ProxyWithContext)

	return nil
}
