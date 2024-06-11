package database

import (
	"fmt"
	"log/slog"
	"time"

	"gorm.io/gorm"
)

type Database struct {
	Session *gorm.DB
}

// MustNewDatabase Establish database connection and migrate tables before
// returning database.Database struct.
func MustNewDatabase(d gorm.Dialector, gc gorm.Config, retryTimes int) Database {
	var (
		DB         *gorm.DB
		err        error
		retryCount int
	)

	slog.Info("Attempting to connect to database")

	err = func() error {
		for i := 0; i <= retryTimes; i++ {
			DB, err = gorm.Open(d, &gc)

			if err == nil {
				retryCount = i
				return nil
			}

			if i == retryTimes {
				return err
			}

			time.Sleep(1 * time.Second)
		}
		return err
	}()
	if err != nil {
		panic(fmt.Sprintf("dataBaseConnect: %v", err))
	}

	slog.Info("Database connection established", "Retry count", retryCount)

	// if err = DB.AutoMigrate(&entity.User{}); err != nil {
	// 	panic(fmt.Sprintf("autoMigrate error: %v", err))
	// }

	return Database{Session: DB}
}
