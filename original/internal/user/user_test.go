package user

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"user-microservice/internal/database/entity"
	"user-microservice/internal/testutil"
	"user-microservice/internal/user/mock"
)

// ━━ TEST SETUP ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

type userSuit struct {
	suite.Suite
	databaseMock *mock.MockdatabaseSession
	service      Service
}

func TestUserSuit(t *testing.T) {
	suite.Run(t, new(userSuit))
}

func (us *userSuit) SetupSuite() {
	DBMock := new(mock.MockdatabaseSession)
	us.databaseMock = DBMock

	us.service = NewService(DBMock)
}

// ━━ TESTS ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

func (us *userSuit) TestList() {
	t := us.T()

	users := []entity.User{
		testutil.NewUser(),
		testutil.NewUser(),
		testutil.NewUser(),
	}

	testCases := map[string]struct {
		mockReturnArgs []any
		expectedReturn any
		expectedError  error
	}{
		"return list of users": {
			[]any{users, nil},
			users,
			nil,
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			us.databaseMock.
				On("ListUsers").
				Return(tc.mockReturnArgs...).
				Once()

			returnUsers, err := us.service.List()
			assert.NoError(t, err, "error listing returnUsers")

			us.databaseMock.AssertCalled(t, "ListUsers")
			us.databaseMock.AssertExpectations(t)

			assert.Equal(
				t,
				tc.expectedReturn,
				returnUsers,
				"expected vs actual returnUsers did not match",
			)
			assert.Equal(t, tc.expectedReturn, err, "expected vs actual err did not match")
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
		expectedReturn interface{}
		expectedError  error
	}{
		"return a user": {
			[]any{int(user.ID)},
			[]any{user, nil},
			int(user.ID),
			user,
			nil,
		},
		"fail to find user": {
			[]any{int(user.ID) + 1},
			[]any{entity.User{}, nil},
			int(user.ID) + 1,
			entity.User{},
			nil,
		},
	}
	for testName, tc := range testCases {
		t.Run(testName, func(t *testing.T) {
			us.databaseMock.
				On("FetchUser", tc.mockInputArgs...).
				Return(tc.mockReturnArgs...).
				Once()

			returnedUser, err := us.service.Fetch(tc.fetchID)

			us.databaseMock.AssertCalled(t, "FetchUser", tc.mockInputArgs...)
			us.databaseMock.AssertExpectations(t)

			assert.Equal(
				t,
				tc.expectedReturn,
				returnedUser,
				"expectedReturn vs actual returnedUser did not match",
			)
			assert.Equal(t, tc.expectedError, err, "expectedReturn vs actual err did not match")
		})
	}
}

func (us *userSuit) TestUpdate() {
	t := us.T()

	user := testutil.NewUser()

	testCases := map[string]struct {
		mockInputArgs  []any
		mockReturnArgs []any
		userID         int
		user           entity.User
		expectedReturn interface{}
		expectedError  error
	}{
		"update user": {
			[]any{int(user.ID), user},
			[]any{user, nil},
			int(user.ID),
			user,
			user,
			nil,
		},
		"fail to update user - ID does not exist": {
			[]any{int(user.ID) + 1, user},
			[]any{entity.User{}, errors.New("test error")},
			int(user.ID) + 1,
			user,
			entity.User{},
			fmt.Errorf("in user.Update: %w", errors.New("test error")),
		},
	}
	for testName, tc := range testCases {
		t.Run(testName, func(t *testing.T) {
			us.databaseMock.
				On("UpdateUser", tc.mockInputArgs...).
				Return(tc.mockReturnArgs...).
				Once()

			returnedUser, err := us.service.Update(tc.userID, tc.user)
			assert.Equal(t, tc.expectedError, err, "expectedReturn vs actual error did not match")

			us.databaseMock.AssertCalled(t, "UpdateUser", tc.mockInputArgs...)
			us.databaseMock.AssertExpectations(t)

			assert.Equal(
				t,
				tc.expectedReturn,
				returnedUser,
				"expectedReturn vs actual returnedUser did not match",
			)
		})
	}
}

func (us *userSuit) TestCreate() {
	t := us.T()

	user := testutil.NewUser()

	testCases := map[string]struct {
		mockInputArgs  []any
		mockReturnArgs []any
		user           entity.User
		expectedReturn interface{}
		expectedError  error
	}{
		"create user": {
			[]any{user},
			[]any{user, nil},
			user,
			int(user.ID),
			nil,
		},
	}
	for testName, tc := range testCases {
		t.Run(testName, func(t *testing.T) {
			us.databaseMock.
				On("CreateUser", tc.mockInputArgs...).
				Return(tc.mockReturnArgs...).
				Once()

			actualReturn, err := us.service.Create(tc.user)
			assert.Equal(t, tc.expectedError, err)

			us.databaseMock.AssertCalled(t, "CreateUser", tc.mockInputArgs...)
			us.databaseMock.AssertExpectations(t)

			assert.Equal(t, tc.expectedReturn, actualReturn)
		})
	}
}

func (us *userSuit) TestDelete() {
	t := us.T()

	user := testutil.NewUser()

	testCases := map[string]struct {
		mockInputArgs  []any
		mockReturnArgs []any
		input          int
		expectedError  error
	}{
		"create user": {
			[]any{int(user.ID)},
			[]any{nil},
			int(user.ID),
			nil,
		},
	}
	for testName, tc := range testCases {
		t.Run(testName, func(t *testing.T) {
			us.databaseMock.
				On("DeleteUser", tc.mockInputArgs...).
				Return(tc.mockReturnArgs...).
				Once()

			err := us.service.Delete(tc.input)
			assert.Equal(t, tc.expectedError, err)

			us.databaseMock.AssertCalled(t, "DeleteUser", tc.mockInputArgs...)
			us.databaseMock.AssertExpectations(t)
		})
	}
}
