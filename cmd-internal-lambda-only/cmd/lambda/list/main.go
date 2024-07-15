package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/jha-captech/user-microservice/internal/config"
	"github.com/jha-captech/user-microservice/internal/database"
	"github.com/jha-captech/user-microservice/internal/handler"
	"github.com/jha-captech/user-microservice/internal/user"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("Startup failed. err: %v", err)
	}
}

func run() error {
	cfg, err := config.NewConfiguration()
	if err != nil {
		return fmt.Errorf("[in run]: %w", err)
	}

	// logger := slog.Default()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel}))

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

	us := user.NewService(db)
	h := handler.NewHandler(logger, us)

	lambda.StartWithOptions(
		h.ListUsersHandler(),
		lambda.WithEnableSIGTERM(func() {
			logger.Info("function container shutting down")
			if err = db.Session.Close(); err != nil {
				logger.Error("error closing database session", "err", err)
			}
		}),
	)

	return nil
}
