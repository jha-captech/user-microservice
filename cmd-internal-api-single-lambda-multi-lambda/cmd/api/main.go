package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httplog/v2"
	"github.com/jha-captech/user-microservice/internal/swagger"

	"github.com/jha-captech/user-microservice/internal/config"
	"github.com/jha-captech/user-microservice/internal/database"
	"github.com/jha-captech/user-microservice/internal/routes"
	"github.com/jha-captech/user-microservice/internal/service"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("Startup failed. err: %v", err)
	}
}

func run() error {
	// Setup
	cfg, err := config.New()
	if err != nil {
		return fmt.Errorf("[in run]: %w", err)
	}

	logger := httplog.NewLogger("user-microservice", httplog.Options{
		LogLevel:        cfg.LogLevel,
		JSON:            false,
		Concise:         true,
		ResponseHeaders: false,
	})

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

	defer func() {
		if err = db.Close(); err != nil {
			logger.Error("Error closing db connection", "err", err)
		}
	}()

	r := chi.NewRouter()

	r.Use(httplog.RequestLogger(logger))
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "PUT", "POST", "DELETE"},
		MaxAge:         300,
	}))

	svs := service.NewUser(db)
	routes.RegisterRoutes(r, logger, svs, routes.WithRegisterHealthRoute(true))

	if cfg.UseSwagger {
		swagger.RunSwagger(r, logger, cfg.HTTP.Domain+cfg.HTTP.Port)
	}

	serverInstance := &http.Server{
		Addr:              cfg.HTTP.Domain + cfg.HTTP.Port,
		IdleTimeout:       time.Minute,
		ReadHeaderTimeout: 500 * time.Millisecond,
		ReadTimeout:       500 * time.Millisecond,
		WriteTimeout:      500 * time.Millisecond,
		Handler:           r,
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

		if err = serverInstance.Shutdown(shutdownCtx); err != nil {
			log.Fatal(fmt.Sprintf("Error shutting down server. err: %v", err))
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
