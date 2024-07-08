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

	"github.com/jha-captech/user-microservice/internal/config"
	"github.com/jha-captech/user-microservice/internal/database"
	"github.com/jha-captech/user-microservice/internal/middleware"
	"github.com/jha-captech/user-microservice/internal/server"
	"github.com/jha-captech/user-microservice/internal/user"
)

func main() {
	run()
}

func run() {
	// Setup
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
	defer db.Session.Close()

	mux := http.NewServeMux()

	stack := middleware.CreateStack(
		middleware.CORSMiddleware(middleware.CORSOptions{
			AllowedOrigins: []string{"*"},
			AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		}),
		middleware.LoggerMiddleware(logger),
		middleware.RecoveryMiddleware(logger),
	)

	us := user.NewService(db)
	h := server.NewHandler(logger, us)
	server.RegisterRoutes(mux, h)

	serverInstance := &http.Server{
		Addr:    cfg.HTTP.Domain + cfg.HTTP.Port,
		Handler: stack(mux),
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

		err := serverInstance.Shutdown(shutdownCtx)
		if err != nil {
			log.Fatal(err)
		}
		serverStopCtx()
	}()

	// Run
	logger.Info(fmt.Sprintf("Server is listening on %s", serverInstance.Addr))
	err := serverInstance.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}

	<-serverCtx.Done()
	logger.Info("Shutdown complete")
}
