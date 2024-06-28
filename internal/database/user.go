package database

import (
	"errors"
	"fmt"
	"log/slog"

	"gorm.io/gorm"

	"user-microservice/internal/database/entity"
)

// ListUsers returns a list of all entity.User objects from the database.
//
// SELECT * FROM "users"
func (db Database) ListUsers() ([]entity.User, error) {
	var users []entity.User
	err := db.Session.Debug().Find(&users).Error
	if err != nil {
		return users, fmt.Errorf("in session.ListUsers: %w", err)
	}

	return users, nil
}

// FetchUser returns am entity.User objects from the database by ID.
//
// SELECT * FROM "users" WHERE ID = $1 ORDER BY "users"."id" LIMIT 1
func (db Database) FetchUser(ID int) (entity.User, error) {
	var user entity.User
	err := db.Session.Debug().Where("ID = ?", ID).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			slog.Info("User not found", "ID", ID)
			return user, nil
		}
		return user, fmt.Errorf("in session.FetchUser: %w", err)
	}

	return user, nil
}

// UpdateUser updates am entity.User objects from the database by ID.
//
// UPDATE "users" SET "first_name"=$1,"last_name"=$2,"role"=$3,"user_id"=$4 WHERE "id" = $5
func (db Database) UpdateUser(ID int, user entity.User) (entity.User, error) {
	// set ID for user to ensure match
	user.ID = uint(ID)

	err := db.Session.Transaction(func(tx *gorm.DB) error {
		if err := tx.Debug().Save(&user).Error; err != nil {
			return fmt.Errorf("in Transaction: %w", err)
		}
		return nil
	})
	if err != nil {
		return entity.User{}, fmt.Errorf("in session.UpdateUser: %w", err)
	}

	return user, nil
}

// CreateUser creates am entity.User objects in the database.
//
// INSERT INTO "users" ("first_name","last_name","role","user_id") VALUES ($1,$2,$3,$4) RETURNING "id"
func (db Database) CreateUser(user entity.User) (entity.User, error) {
	// set ID to 0 so that it is auto generated
	user.ID = uint(0)

	err := db.Session.Transaction(func(tx *gorm.DB) error {
		if err := tx.Debug().Create(&user).Error; err != nil {
			return fmt.Errorf("in Transaction: %w", err)
		}
		return nil
	})
	if err != nil {
		return entity.User{}, fmt.Errorf("in session.CreateUser: %w", err)
	}

	return user, nil
}

// DeleteUser deletes am entity.User objects from the database by ID.
//
// DELETE FROM "users" WHERE "users"."id" = $1
func (db Database) DeleteUser(ID int) error {
	if err := db.Session.Debug().Delete(&entity.User{}, ID).Error; err != nil {
		return fmt.Errorf("in session.DeleteUser: %w", err)
	}

	return nil
}
