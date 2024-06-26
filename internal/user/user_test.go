package user

import (
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"user-microservice/internal/database/entity"
	"user-microservice/internal/testutil"
)

// MOCKS

type databaseMock struct {
	mock.Mock
}

func (dsm *databaseMock) ListUsers() ([]entity.User, error) {
	args := dsm.Called()
	return args.Get(0).([]entity.User), args.Error(1)
}

func (dsm *databaseMock) FetchUser(ID int) (entity.User, error) {
	args := dsm.Called(ID)
	return args.Get(0).(entity.User), args.Error(1)
}

// TEST SETUP

type userSuit struct {
	suite.Suite
	databaseMock *databaseMock
	service      Service
}

func TestUserSuit(t *testing.T) {
	suite.Run(t, new(userSuit))
}

func (us *userSuit) SetupSuite() {
	DBMock := new(databaseMock)
	us.databaseMock = DBMock

	logger := slog.Default()

	us.service = NewService(DBMock, logger)
}

// TESTS

func (us *userSuit) TestList() {
	t := us.T()

	users := []entity.User{
		testutil.NewUser(),
		testutil.NewUser(),
		testutil.NewUser(),
	}

	testCases := map[string]struct {
		mockReturnArgs []any
		expected       any
	}{
		"return list of users": {
			[]any{users, nil},
			users,
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

			us.databaseMock.AssertCalled(t, "ListUsers")
			us.databaseMock.AssertExpectations(t)

			assert.Equal(t, users, tc.expected, "expected vs actual users did not match")
		})
	}
}

func (us *userSuit) TestFetch() {
	t := us.T()

	user := testutil.NewUser()

	testCases := map[string]struct {
		mockInputArgs  []any
		mockReturnArgs []any
		fetchID        int
		expected       interface{}
	}{
		"return a user": {
			[]any{int(user.ID)},
			[]any{user, nil},
			int(user.ID),
			user,
		},
		"fail to find user": {
			[]any{int(user.ID) + 1},
			[]any{entity.User{}, nil},
			int(user.ID) + 1,
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
