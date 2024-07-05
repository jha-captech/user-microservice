package testutil

import (
	"github.com/go-faker/faker/v4"

	"user-microservice/internal/database/entity"
)

type Options[T any] func(*T)

func WithID(ID int) Options[entity.User] {
	return func(user *entity.User) {
		user.ID = uint(ID)
	}
}

func WithFirstName(firstName string) Options[entity.User] {
	return func(user *entity.User) {
		user.FirstName = firstName
	}
}

func WithLastName(lastName string) Options[entity.User] {
	return func(user *entity.User) {
		user.LastName = lastName
	}
}

func WithRole(role string) Options[entity.User] {
	return func(user *entity.User) {
		user.Role = role
	}
}

func WithUserID(userID int) Options[entity.User] {
	return func(user *entity.User) {
		user.UserID = uint(userID)
	}
}

func NewUser(options ...Options[entity.User]) entity.User {
	user := entity.User{
		ID:        uint(faker.RandomUnixTime()),
		FirstName: faker.FirstName(),
		LastName:  faker.LastName(),
		Role:      faker.Word(),
		UserID:    uint(faker.RandomUnixTime()),
	}

	for _, option := range options {
		option(&user)
	}

	return user
}

type newUsersSettings struct {
	startID int
}

func WithIDStartRange(ID int) Options[newUsersSettings] {
	return func(settings *newUsersSettings) {
		settings.startID = ID
	}
}

func NewUsers(count int, options ...Options[newUsersSettings]) []entity.User {
	settings := newUsersSettings{
		startID: -1,
	}

	for _, option := range options {
		option(&settings)
	}

	users := make([]entity.User, count)
	for i := 0; i < count; i++ {
		var opts []Options[entity.User]
		if settings.startID > 0 {
			opts = append(opts, WithID(settings.startID))
			settings.startID++
		}
		user := NewUser(opts...)
		users[i] = user
	}

	return users
}
