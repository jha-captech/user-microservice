package user

import (
	"fmt"

	"user-microservice/internal/database/entity"
)

type databaseSession interface {
	ListUsers() ([]entity.User, error)
	FetchUser(ID int) (entity.User, error)
}

type Service struct {
	Database databaseSession
}

func NewService(db databaseSession) Service {
	return Service{
		Database: db,
	}
}

func (s Service) List() ([]entity.User, error) {
	users, err := s.Database.ListUsers()
	if err != nil {
		return []entity.User{}, fmt.Errorf("in user.List: %w", err)
	}
	return users, nil
}

func (s Service) Fetch(ID int) (entity.User, error) {
	user, err := s.Database.FetchUser(ID)
	if err != nil {
		return entity.User{}, fmt.Errorf("in user.Fetch: %w", err)
	}
	return user, nil
}
