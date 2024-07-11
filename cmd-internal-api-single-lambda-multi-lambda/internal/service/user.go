package service

import (
	"database/sql"
	"fmt"

	"github.com/jha-captech/user-microservice/internal/models"
)

type Service struct {
	Database *sql.DB
}

// NewService returns a new Service struct.
func NewService(db *sql.DB) Service {
	return Service{
		Database: db,
	}
}

// ListUsers returns a list of all User objects from the Database.
func (s Service) ListUsers() ([]models.User, error) {
	rows, err := s.Database.Query(
		`
		SELECT 
		    * 
		FROM
		    "users" 
		`,
	)
	if err != nil {
		return []models.User{}, fmt.Errorf("[in ListUsers]:, %w", err)
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err = rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Role, &user.UserID)
		if err != nil {
			return []models.User{}, fmt.Errorf("[in ListUsers]:, %w", err)
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return []models.User{}, fmt.Errorf("[in ListUsers]:, %w", err)
	}

	return users, nil
}

// FetchUser returns am User objects from the Database by ID.
func (s Service) FetchUser(ID int) (models.User, error) {
	var user models.User
	err := s.Database.
		QueryRow(
			`
			SELECT
				*
			FROM
				"users"
			WHERE
				id = $1
			ORDER BY
				"users"."id"
			LIMIT 1
			`,
			ID,
		).
		Scan(&user.ID, &user.FirstName, &user.LastName, &user.Role, &user.UserID)
	if err != nil {
		return models.User{}, fmt.Errorf("[in FetchUser]:, %w", err)
	}

	return user, nil
}

// UpdateUser updates am User objects from the Database by ID.
func (s Service) UpdateUser(ID int, user models.User) (models.User, error) {
	_, err := s.Database.Exec(
		`
		UPDATE
			"users"
		SET
			"first_name" = $1,
			"last_name" = $2,
			"role" = $3,
			"user_id" = $4
		WHERE
			"id" = $5
		`,
		user.FirstName,
		user.LastName,
		user.Role,
		user.UserID,
		ID,
	)
	if err != nil {
		return models.User{}, fmt.Errorf("[in UpdateUser]: %w", err)
	}

	user.ID = uint(ID)
	return user, nil
}

// CreateUser creates am User objects in the Database.
func (s Service) CreateUser(user models.User) (int, error) {
	var ID int
	err := s.Database.QueryRow(
		`
		INSERT INTO "users" ("first_name", "last_name", "role", "user_id")
			VALUES ($1, $2, $3, $4)
		RETURNING "id"
		`,
		user.FirstName,
		user.LastName,
		user.Role,
		user.UserID,
	).Scan(&ID)
	if err != nil {
		return 0, fmt.Errorf("[in CreateUser]: %w", err)
	}

	return ID, nil
}

// DeleteUser deletes am User objects from the Database by ID.
func (s Service) DeleteUser(ID int) error {
	_, err := s.Database.Exec(
		`
		DELETE FROM "users"
		WHERE "users"."id" = $1
		`,
		ID,
	)
	if err != nil {
		return fmt.Errorf("[in DeleteUser]: %w", err)
	}

	return nil
}
