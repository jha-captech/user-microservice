package main

import (
	"database/sql"
	"errors"
	"fmt"
)

type UserService struct {
	DB Database
}

// NewUserService returns a new UserService struct.
func NewUserService(db Database) UserService {
	return UserService{
		DB: db,
	}
}

// ListUsers returns a list of all User objects from the Database.
func (us UserService) ListUsers() ([]User, error) {
	rows, err := us.DB.Session.Query(
		`
		SELECT 
		    * 
		FROM
		    "users" 
		`,
	)
	if err != nil {
		return []User{}, fmt.Errorf("in ListUsers:, %w", err)
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err = rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Role, &user.UserID)
		if err != nil {
			return []User{}, fmt.Errorf("in ListUsers:, %w", err)
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return []User{}, fmt.Errorf("in ListUsers:, %w", err)
	}

	return users, nil
}

// FetchUser returns am User objects from the Database by ID.
func (us UserService) FetchUser(ID int) (User, error) {
	var user User
	err := us.DB.Session.
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
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, nil
		}
		return User{}, fmt.Errorf("in ListUsers:, %w", err)
	}

	return user, nil
}

// UpdateUser updates am User objects from the Database by ID.
func (us UserService) UpdateUser(ID int, user User) (User, error) {
	_, err := us.DB.Session.Exec(
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
		return User{}, fmt.Errorf("in UpdateUser: %w", err)
	}

	user.ID = uint(ID)
	return user, nil
}

// CreateUser creates am User objects in the Database.
func (us UserService) CreateUser(user User) (int, error) {
	var ID int
	err := us.DB.Session.QueryRow(
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
		return 0, fmt.Errorf("in CreateUser: %w", err)
	}

	return ID, nil
}

// DeleteUser deletes am User objects from the Database by ID.
func (us UserService) DeleteUser(ID int) error {
	_, err := us.DB.Session.Exec(
		`
		DELETE FROM "users"
		WHERE "users"."id" = $1
		`,
		ID,
	)
	if err != nil {
		return fmt.Errorf("could not delete user: %v", err)
	}

	return nil
}
