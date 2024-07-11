package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/go-chi/chi/v5"
	"github.com/jha-captech/user-microservice/internal/database"
	"github.com/jha-captech/user-microservice/internal/handlers"
	"github.com/jha-captech/user-microservice/internal/middleware"
	"github.com/jha-captech/user-microservice/internal/routes"
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
	HTTP struct {
		Domain              string `env:"HTTP_DOMAIN,required"`
		Port                string `env:"HTTP_PORT,required"`
		ShutdownGracePeriod int    `env:"HTTP_SHUTDOWN_GRACE_PERIOD,required"`
	}
}

func main() {
	if err := run(); err != nil {
		log.Fatalf("Startup failed. err: %v", err)
	}
}

func run() error {
	// Setup
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

	stack := middleware.CreateStack(
		middleware.CORSMiddleware(middleware.CORSOptions{
			AllowedOrigins: []string{"*"},
			AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		}),
		middleware.LoggerMiddleware(logger),
		middleware.RecoveryMiddleware(logger),
	)

	us := user.NewService(db)
	h := handlers.New(logger, us)
	routes.RegisterRoutes(r, h)

	serverInstance := &http.Server{
		Addr:    cfg.HTTP.Domain + cfg.HTTP.Port,
		Handler: stack(r),
	}

	// Graceful shutdown
	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig

		fmt.Println()
		logger.Info("Shutdown signal received")

		shutdownCtx, _ := context.WithTimeout(
			serverCtx, time.Duration(cfg.HTTP.ShutdownGracePeriod)*time.Second,
		)

		go func() {
			<-shutdownCtx.Done()
			if errors.Is(shutdownCtx.Err(), context.DeadlineExceeded) {
				log.Fatal("graceful shutdown timed out.. forcing exit.")
			}
		}()

		err = serverInstance.Shutdown(shutdownCtx)
		if err != nil {
			log.Fatal(err)
		}
		serverStopCtx()
	}()

	// Run
	logger.Info(fmt.Sprintf("Server is listening on %s", serverInstance.Addr))
	err = serverInstance.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	<-serverCtx.Done()
	logger.Info("Shutdown complete")
	return nil
}