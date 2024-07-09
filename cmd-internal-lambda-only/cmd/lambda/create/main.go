package main

import (
	"fmt"
	"log/slog"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/jha-captech/user-microservice/internal/config"
	"github.com/jha-captech/user-microservice/internal/database"
	"github.com/jha-captech/user-microservice/internal/handler"
	"github.com/jha-captech/user-microservice/internal/user"
)

func main() {
	cfg := config.MustNewConfiguration()

	logger := slog.Default()

	db := database.MustNewDatabase(
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

	us := user.NewService(db)
	h := handler.NewHandler(logger, us)

	lambda.StartWithOptions(
		h.CreateUsersHandler(),
		lambda.WithEnableSIGTERM(func() {
			logger.Info("function container shutting down")
			if err := db.Session.Close(); err != nil {
				logger.Error("error closing database session", "err", err)
			}
		}),
	)
}
