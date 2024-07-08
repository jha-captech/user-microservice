package main

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
	"github.com/jha-captech/user-microservice/internal/config"
	"github.com/jha-captech/user-microservice/internal/database"
	"github.com/jha-captech/user-microservice/internal/server"
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
	cfg := config.MustNewConfiguration[configuration]()

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

	mux := http.NewServeMux()

	us := user.NewService(db)
	h := server.NewHandler(logger, us)
	server.RegisterRoutes(mux, h, server.WithEnableHealthCheck(false))

	lambda.StartWithOptions(
		httpadapter.New(mux).ProxyWithContext,
		lambda.WithEnableSIGTERM(func() {
			logger.Info("function container shutting down")
			if err := db.Session.Close(); err != nil {
				logger.Error("error closing database session", "err", err)
			}
		}),
	)
}
