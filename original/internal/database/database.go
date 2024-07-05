package database

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"gorm.io/gorm"

	"user-microservice/internal/database/entity"
)

type Database struct {
	Session *gorm.DB
}

type Options func(*dbSetupOptions)

type dbSetupOptions struct {
	connectionRetry int
	runMigrations   bool
	gormConfig      gorm.Config
	logger          *slog.Logger
}

func WithRetryCount(count int) Options {
	return func(setup *dbSetupOptions) {
		setup.connectionRetry = count
	}
}

func WithAutoMigrate(runMigration bool) Options {
	return func(setup *dbSetupOptions) {
		setup.runMigrations = runMigration
	}
}

func WithGormConfig(config gorm.Config) Options {
	return func(setup *dbSetupOptions) {
		setup.gormConfig = config
	}
}

func WithLogger(logger *slog.Logger) Options {
	return func(setup *dbSetupOptions) {
		setup.logger = logger
	}
}

// MustNewDatabase Establish session connection and migrate tables before
// returning database.Database struct.
func MustNewDatabase(d gorm.Dialector, options ...Options) Database {
	opts := dbSetupOptions{
		connectionRetry: 5,
		runMigrations:   false,
		gormConfig:      gorm.Config{},
		logger:          slog.New(slog.NewJSONHandler(os.Stdout, nil)),
	}

	for _, option := range options {
		option(&opts)
	}

	var (
		DB         *gorm.DB
		err        error
		retryCount int
	)

	slog.Info("Attempting to connect to session")

	err = func() error {
		for i := 0; i <= opts.connectionRetry; i++ {
			DB, err = gorm.Open(d, &opts.gormConfig)

			if err == nil {
				retryCount = i
				return nil
			}

			if i == opts.connectionRetry {
				return err
			}

			time.Sleep(1 * time.Second)
		}
		return err
	}()
	if err != nil {
		panic(fmt.Sprintf("dataBaseConnect: %v", err))
	}

	opts.logger.Info("Database connection established", "Retry count", retryCount)

	if opts.runMigrations {
		if err = DB.AutoMigrate(&entity.User{}); err != nil {
			panic(fmt.Sprintf("autoMigrate error: %v", err))
		}
		opts.logger.Info("Database migration successful")
	}

	return Database{Session: DB}
}
