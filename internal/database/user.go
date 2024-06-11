package database

import (
	"errors"
	"fmt"
	"log/slog"

	"gorm.io/gorm"

	"user-microservice/internal/database/entity"
)

func (db Database) ListUsers() ([]entity.User, error) {
	var users []entity.User
	err := db.Session.Find(&users).Error
	if err != nil {
		return users, fmt.Errorf("in database.ListUsers: %w", err)
	}

	return users, nil
}

func (db Database) FetchUser(ID int) (entity.User, error) {
	var user entity.User
	err := db.Session.Where("ID = ?", ID).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			slog.Info("User not found", "ID", ID)
			return user, nil
		}
		return user, fmt.Errorf("in database.FetchUser: %w", err)
	}

	return user, nil
}
