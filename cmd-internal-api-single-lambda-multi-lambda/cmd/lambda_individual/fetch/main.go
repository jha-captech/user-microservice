package main

import (
	"fmt"
	"log"
	"log/slog"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
	"github.com/go-chi/chi/v5"
	"github.com/jha-captech/user-microservice/internal/config"
	"github.com/jha-captech/user-microservice/internal/database"
	"github.com/jha-captech/user-microservice/internal/handlers"
	"github.com/jha-captech/user-microservice/internal/service"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("Startup failed. err: %v", err)
	}
}

func run() error {
	cfg, err := config.New()
	if err != nil {
		return fmt.Errorf("[in run]: %w", err)
	}
	logger := slog.Default()

	db, err := database.New(
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

	defer db.Close()

	r := chi.NewRouter()

	us := service.NewService(db)
	h := handlers.New(logger, us)

	r.Get("/api/user/{ID}", h.HandleFetchUser())

	lambda.Start(httpadapter.New(r).ProxyWithContext)

	return nil
}
