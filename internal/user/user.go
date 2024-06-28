package user

import (
	"fmt"

	"user-microservice/internal/database/entity"
)

type databaseSession interface {
	ListUsers() ([]entity.User, error)
	FetchUser(int) (entity.User, error)
	UpdateUser(int, entity.User) (entity.User, error)
	CreateUser(entity.User) (entity.User, error)
	DeleteUser(int) error
}

type Service struct {
	Database databaseSession
}

// NewService returns a new instance of the Service struct.
func NewService(db databaseSession) Service {
	return Service{
		Database: db,
	}
}

// List returns a list of type []entity.User.
func (s Service) List() ([]entity.User, error) {
	users, err := s.Database.ListUsers()
	if err != nil {
		return []entity.User{}, fmt.Errorf("in user.List: %w", err)
	}
	return users, nil
}

// Fetch returns an object of type entity.User.
func (s Service) Fetch(ID int) (entity.User, error) {
	user, err := s.Database.FetchUser(ID)
	if err != nil {
		return entity.User{}, fmt.Errorf("in user.Fetch: %w", err)
	}
	return user, nil
}

// Update updates a entity.User object bu ID.
func (s Service) Update(ID int, user entity.User) (entity.User, error) {
	user, err := s.Database.UpdateUser(ID, user)
	if err != nil {
		return entity.User{}, fmt.Errorf("in user.Update: %w", err)
	}
	return user, nil
}

// Create creates an entity.User object based on a entity.User passed in.
func (s Service) Create(user entity.User) (int, error) {
	user, err := s.Database.CreateUser(user)
	if err != nil {
		return 0, fmt.Errorf("in user.Create: %w", err)
	}
	return int(user.ID), nil
}

// Delete deletes a entity.User object by ID.
func (s Service) Delete(ID int) error {
	if err := s.Database.DeleteUser(ID); err != nil {
		return fmt.Errorf("in user.Delete: %w", err)
	}
	return nil
}
