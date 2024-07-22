package service

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jha-captech/user-microservice/internal/models"
)

type User struct {
	database *sql.DB
}

// NewUser returns a new User struct.
func NewUser(db *sql.DB) *User {
	return &User{
		database: db,
	}
}

// ListUsers returns a list of all User objects from the database.
func (s User) ListUsers(ctx context.Context) ([]models.User, error) {
	rows, err := s.database.QueryContext(
		ctx,
		`SELECT * FROM "users"`,
	)
	if err != nil {
		return []models.User{}, fmt.Errorf("[in ListUsers]: %w", err)
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err = rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Role, &user.UserID)
		if err != nil {
			return []models.User{}, fmt.Errorf("[in ListUsers]: %w", err)
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return []models.User{}, fmt.Errorf("[in ListUsers]: %w", err)
	}

	return users, nil
}

// FetchUser returns am User objects from the database by ID.
func (s User) FetchUser(ctx context.Context, ID int) (models.User, error) {
	var user models.User
	err := s.database.
		QueryRowContext(
			ctx,
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
		return models.User{}, fmt.Errorf("[in FetchUser]: %w", err)
	}

	return user, nil
}

// UpdateUser updates am User objects from the database by ID.
func (s User) UpdateUser(ctx context.Context, ID int, user models.User) (models.User, error) {
	_, err := s.database.ExecContext(
		ctx,
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

// CreateUser creates am User objects in the database.
func (s User) CreateUser(ctx context.Context, user models.User) (int, error) {
	var ID int
	err := s.database.QueryRowContext(
		ctx,
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

// DeleteUser deletes am User objects from the database by ID.
func (s User) DeleteUser(ctx context.Context, ID int) error {
	_, err := s.database.ExecContext(
		ctx,
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
