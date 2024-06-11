package user

import (
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"user-microservice/internal/database/entity"
)

// mocks

type databaseSessionMock struct {
	Users []entity.User
}

func newDatabaseSessionMock(userCount int) databaseSessionMock {
	var mockUsers []entity.User
	for i := 0; i < userCount; i++ {
		mockUser := entity.User{
			ID:        uint(faker.RandomUnixTime()),
			FirstName: faker.FirstName(),
			LastName:  faker.LastName(),
			Role:      faker.Word(),
			UserID:    uint(faker.RandomUnixTime()),
		}
		mockUsers = append(mockUsers, mockUser)
	}
	return databaseSessionMock{Users: mockUsers}
}

func (dsm databaseSessionMock) ListUsers() ([]entity.User, error) {
	return dsm.Users, nil
}

func (dsm databaseSessionMock) FetchUser(ID int) (entity.User, error) {
	for _, mockUser := range dsm.Users {
		if mockUser.ID != uint(ID) {
			continue
		}
		return mockUser, nil
	}
	return entity.User{}, nil
}

// test setup

type userSuit struct {
	suite.Suite
	userDataMock    []entity.User
	userServiceMock Service
}

func TestUserSuit(t *testing.T) {
	suite.Run(t, new(userSuit))
}

func (us *userSuit) SetupTest() {
	databaseMock := newDatabaseSessionMock(5)
	us.userDataMock = databaseMock.Users

	us.userServiceMock = NewService(databaseMock)
}

// tests

func (us *userSuit) TestList() {
	t := us.T()

	testCases := []struct {
		name     string
		expected interface{}
	}{
		{
			"return list of users",
			us.userDataMock,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			users, err := us.userServiceMock.List()
			assert.NoError(t, err, "error listing users")
			assert.Equal(t, users, tc.expected, "expected vs actual users did not match")
		})
	}
}

func (us *userSuit) TestFetch() {
	t := us.T()

	testCases := map[string]struct {
		ID       int
		expected interface{}
	}{
		"return a user": {
			int(us.userDataMock[1].ID),
			us.userDataMock[1],
		},
		"also return a user": {
			2,
			func(ID int) entity.User {
				for _, user := range us.userDataMock {
					if user.ID == uint(ID) {
						return user
					}
				}
				return entity.User{}
			}(2),
		},
	}
	for testName, tc := range testCases {
		t.Run(testName, func(t *testing.T) {
			user, err := us.userServiceMock.Fetch(tc.ID)
			assert.NoError(t, err, "error listing user")
			assert.Equal(t, user, tc.expected, "expected vs actual user did not match")
		})
	}
}
