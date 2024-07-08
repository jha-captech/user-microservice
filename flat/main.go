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
)

func main() {
	run()
}

func run() {
	// Setup
	config := MustNewConfiguration()

	logger := slog.Default()

	db := MustNewDatabase(
		fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
			config.Database.Host,
			config.Database.User,
			config.Database.Password,
			config.Database.Name,
			config.Database.Port,
		),
		logger,
		config.Database.ConnectionRetry,
	)
	defer db.Session.Close()

	us := NewUserService(db)

	h := newHandler(logger, us)

	mux := http.NewServeMux()

	stack := CreateStack(
		CORSMiddleware(CORSOptions{
			allowedOrigins: []string{"*"},
			allowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		}),
		LoggerMiddleware(logger),
		RecoveryMiddleware(logger),
	)

	RegisterRoutes(mux, h)

	server := &http.Server{
		Addr:    config.HTTP.Domain + config.HTTP.Port,
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
			serverCtx, time.Duration(config.HTTP.ShutdownGracePeriod)*time.Second,
		)

		go func() {
			<-shutdownCtx.Done()
			if errors.Is(shutdownCtx.Err(), context.DeadlineExceeded) {
				log.Fatal("graceful shutdown timed out.. forcing exit.")
			}
		}()

		err := server.Shutdown(shutdownCtx)
		if err != nil {
			log.Fatal(err)
		}
		serverStopCtx()
	}()

	// Run
	logger.Info(fmt.Sprintf("Server is listening on %s", server.Addr))
	err := server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}

	<-serverCtx.Done()
	logger.Info("Shutdown complete")
}
