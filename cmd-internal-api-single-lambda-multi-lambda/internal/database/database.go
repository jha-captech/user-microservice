package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"time"

	_ "github.com/lib/pq"
)

type Database struct {
	Session *sql.DB
}

// NewDatabase Establish Session connection and migrate tables before
// returning Database.Database struct.
func NewDatabase(connectionString string, logger *slog.Logger, retryCount int) (Database, error) {
	logger.Info("Attempting to connect to Database")
	db, err := retryWithReturn(retryCount, 100*time.Millisecond, func() (*sql.DB, error) {
		return sql.Open("postgres", connectionString)
	})
	if err != nil {
		return Database{}, fmt.Errorf(
			"[in NewDatabase] Failed to connect to Database after %d attempts: %w",
			retryCount,
			err,
		)
	}

	logger.Info("Attempting to ping Database")
	err = retry(retryCount, 500*time.Millisecond, func() error {
		return db.Ping()
	})
	if err != nil {
		db.Close()
		return Database{}, fmt.Errorf(
			"[in NewDatabase] Failed to ping Database after %d attempts: %w",
			retryCount,
			err,
		)
	}

	logger.Info("Database connection established")

	return Database{Session: db}, nil
}

// retry will retry a given function n times with a wait of a given duration between each retry attempt.
func retry(retryCount int, waitTime time.Duration, fn func() error) error {
	_, err := retryWithReturn(retryCount, waitTime, func() (any, error) {
		return nil, fn()
	})
	return err
}

// retryWithReturn will retry a given function n times with a wait of a given duration between each
// retry attempt. retryWithReturn is intended for functions where a return values is needed.
func retryWithReturn[T any](retryCount int, waitTime time.Duration, fn func() (T, error)) (T, error) {
	if retryCount < 1 {
		return *new(T), errors.New("retryCount of less than 1 is not permitted")
	}
	for i := 0; i < retryCount; i++ {
		t, err := fn()
		if err != nil {
			if i == retryCount-1 {
				return t, err
			}
			time.Sleep(waitTime)
			continue
		}
		return t, nil
	}
	return *new(T), errors.New("default return reached in retryWithReturn")
}
