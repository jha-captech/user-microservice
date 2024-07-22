package database

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/jha-captech/user-microservice/internal/log"
	_ "github.com/lib/pq"
)

// New Establish Session connection and migrate tables before
// returning database.database struct.
func New(ctx context.Context, connectionString string, logger log.Logger, retryDuration time.Duration) (*sql.DB, error) {
	logger.Info("Attempting to connect to database")
	retryCount := 0
	db, err := retryResult(ctx, retryDuration, func() (*sql.DB, error) {
		retryCount++
		return sql.Open("postgres", connectionString)
	})
	if err != nil {
		return nil, fmt.Errorf(
			"[in New] Failed to connect to database with retry duration of %s and %d attempts: %w",
			retryDuration,
			retryCount,
			err,
		)
	}

	logger.Info("Successfully connected to database", "retry count", retryCount)

	logger.Info("Attempting to ping database")

	retryCount = 0
	err = retry(ctx, retryDuration, func() error {
		retryCount++
		return db.Ping()
	})
	if err != nil {
		db.Close()
		return nil, fmt.Errorf(
			"[in New] Failed to ping database with retry duration of %s and %d attempts: %w",
			retryDuration,
			retryCount,
			err,
		)
	}
	logger.Info("Successfully pinged database", "retry count", retryCount)

	logger.Info("database connection established")

	return db, nil
}

func retry(ctx context.Context, maxDuration time.Duration, retryFunc func() error) error {
	_, err := retryResult(ctx, maxDuration, func() (any, error) {
		return nil, retryFunc()
	})
	return err
}

func retryResult[T any](ctx context.Context, maxDuration time.Duration, retryFunc func() (T, error)) (T, error) {
	var (
		returnData T
		err        error
	)
	const maxBackoffMilliseconds = 2_000.0

	ctx, cancelFunc := context.WithTimeout(ctx, maxDuration)
	defer cancelFunc()

	go func() {
		counter := 1.0
		for {
			counter++
			returnData, err = retryFunc()
			if err != nil {
				waitMilliseconds := math.Min(
					math.Pow(counter, 2)+float64(rand.Intn(10)),
					maxBackoffMilliseconds,
				)
				time.Sleep(time.Duration(waitMilliseconds) * time.Millisecond)
				continue
			}
			cancelFunc()
			return
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return returnData, err
		}
	}
}
