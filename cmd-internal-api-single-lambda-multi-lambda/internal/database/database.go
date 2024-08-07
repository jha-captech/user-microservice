package database

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

type sLogger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}

// New Establish Session connection and migrate tables before
// returning database.database struct.
func New(connectionString string, logger sLogger, retryCount int) (*sql.DB, error) {
	logger.Info("Attempting to connect to database")
	db, err := retryWithReturn(retryCount, 100*time.Millisecond, func() (*sql.DB, error) {
		return sql.Open("postgres", connectionString)
	})
	if err != nil {
		return nil, fmt.Errorf(
			"[in New] Failed to connect to database after %d attempts: %w",
			retryCount,
			err,
		)
	}

	logger.Info("Attempting to ping database")
	err = retry(retryCount, 500*time.Millisecond, func() error {
		return db.Ping()
	})
	if err != nil {
		db.Close()
		return nil, fmt.Errorf(
			"[in New] Failed to ping database after %d attempts: %w",
			retryCount,
			err,
		)
	}

	logger.Info("database connection established")

	return db, nil
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
