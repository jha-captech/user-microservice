package user

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"user-microservice/internal/database/entity"
)

// mocks

type databaseSessionMock struct {
	mock.Mock
}

func (dsm *databaseSessionMock) ListUsers() ([]entity.User, error) {
	args := dsm.Called()
	return args.Get(0).([]entity.User), args.Error(1)
}

func (dsm *databaseSessionMock) FetchUser(ID int) (entity.User, error) {
	args := dsm.Called(ID)
	return args.Get(0).(entity.User), args.Error(1)
}

// test setup

type userSuit struct {
	suite.Suite
	databaseMock *databaseSessionMock
	service      Service
}

func TestUserSuit(t *testing.T) {
	suite.Run(t, new(userSuit))
}

func (us *userSuit) SetupTest() {
	databaseMock := new(databaseSessionMock)
	us.databaseMock = databaseMock

	us.service = NewService(databaseMock)
}

// tests

func (us *userSuit) TestList() {
	t := us.T()

	testCases := map[string]struct {
		mockReturnArgs []any
		expected       any
	}{
		"return list of users": {
			[]any{[]entity.User{}, nil},
			[]entity.User{},
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			us.databaseMock.
				On("ListUsers").
				Return(tc.mockReturnArgs...).
				Once()

			users, err := us.service.List()
			assert.NoError(t, err, "error listing users")

			us.databaseMock.AssertCalled(t, "FetchUser")
			us.databaseMock.AssertExpectations(t)

			assert.Equal(t, users, tc.expected, "expected vs actual users did not match")
		})
	}
}

func (us *userSuit) TestFetch() {
	t := us.T()

	testCases := map[string]struct {
		mockInputArgs  []any
		mockReturnArgs []any
		fetchID        int
		expected       interface{}
	}{
		"return a user": {
			[]any{1},
			[]any{entity.User{}, nil},
			1,
			entity.User{},
		},
	}
	for testName, tc := range testCases {
		t.Run(testName, func(t *testing.T) {
			us.databaseMock.
				On("FetchUser", tc.mockInputArgs...).
				Return(tc.mockReturnArgs...).
				Once()

			user, err := us.service.Fetch(tc.fetchID)
			assert.NoError(t, err, "error listing user")

			us.databaseMock.AssertCalled(t, "FetchUser", tc.mockInputArgs...)
			us.databaseMock.AssertExpectations(t)

			assert.Equal(t, user, tc.expected, "expected vs actual user did not match")
		})
	}
}
