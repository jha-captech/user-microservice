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

	"user-microservice/cmd/http/route"
	"user-microservice/internal/database"
	"user-microservice/internal/user"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"gorm.io/driver/postgres"
)

func main() {
	logger, config, router := SetupServer()
	fmt.Printf("%+v", config)

	server := &http.Server{
		Addr:    config.HTTP.Domain + config.HTTP.Port,
		Handler: router,
	}

	// Server run context
	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	// Listen for syscall signals for process to interrupt/quit
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig

		fmt.Println()
		logger.Info("Shutdown signal received")

		// Shutdown signal with grace period of 30 seconds
		shutdownCtx, _ := context.WithTimeout(
			serverCtx, time.Duration(config.HTTP.ShutdownGracePeriod)*time.Second,
		)

		go func() {
			<-shutdownCtx.Done()
			if errors.Is(shutdownCtx.Err(), context.DeadlineExceeded) {
				log.Fatal("graceful shutdown timed out.. forcing exit.")
			}
		}()

		// Trigger graceful shutdown
		err := server.Shutdown(shutdownCtx)
		if err != nil {
			log.Fatal(err)
		}
		serverStopCtx()
	}()

	// Run the server
	logger.Info(fmt.Sprintf("Server is listening on %s", server.Addr))
	err := server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}

	// Wait for server context to be stopped
	<-serverCtx.Done()
	logger.Info("Shutdown complete")
}

func SetupServer() (*slog.Logger, configuration, *chi.Mux) {
	config := mustNewConfiguration()

	logger := newLogger(false)

	db := database.MustNewDatabase(
		postgres.Open(
			fmt.Sprintf(
				"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
				config.Database.Host,
				config.Database.User,
				config.Database.Password,
				config.Database.Name,
				config.Database.Port,
			),
		),
		database.WithLogger(logger),
		database.WithRetryCount(5),
		database.WithAutoMigrate(true),
	)

	us := user.NewService(db)

	h := route.NewHandler(us, logger)

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	route.SetUpRoutes(r, h)

	return logger, config, r
}
