package route

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/go-faker/faker/v4"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"user-microservice/internal/database/entity"
)

// mocks

type serviceMock struct {
	mock.Mock
}

func (sm *serviceMock) List() ([]entity.User, error) {
	args := sm.Called()
	return args.Get(0).([]entity.User), args.Error(1)
}

func (sm *serviceMock) Fetch(ID int) (entity.User, error) {
	args := sm.Called(ID)
	return args.Get(0).(entity.User), args.Error(1)
}

// test setup

func generateUsers(count int) []entity.User {
	var users []entity.User
	for i := 0; i < count; i++ {
		user := entity.User{
			ID:        uint(faker.RandomUnixTime()),
			FirstName: faker.FirstName(),
			LastName:  faker.LastName(),
			Role:      faker.Word(),
			UserID:    uint(faker.RandomUnixTime()),
		}
		users = append(users, user)
	}
	return users
}

type routerSuit struct {
	suite.Suite
	router   *gin.Engine
	userMock *serviceMock
	users    []entity.User
}

func TestRouterSuit(t *testing.T) {
	suite.Run(t, new(routerSuit))
}

func (rs *routerSuit) SetupTest() {
	rs.users = generateUsers(5)

	userServiceMock := new(serviceMock)
	rs.userMock = userServiceMock
	handler := NewHandler(userServiceMock)

	rs.router = gin.Default()
	SetUpRoutes(rs.router, handler)
}

// tests

func (rs *routerSuit) TestHealthCheck() {
	t := rs.T()

	testCases := map[string]struct {
		method         string
		path           string
		expectedStatus int
	}{
		"200": {
			http.MethodGet,
			"/health-check/",
			http.StatusOK,
		},
		"404": {
			http.MethodPut,
			"/health-check/",
			http.StatusNotFound,
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(tc.method, tc.path, nil)
			rs.router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatus, w.Code, "Wrong code received")
		})
	}
}

func (rs *routerSuit) TestUserList() {
	t := rs.T()

	testCases := map[string]struct {
		mockReturnArgs []any
		method         string
		path           string
		expectedStatus int
		expectedBody   any
	}{
		"200 - Good call": {
			[]any{rs.users, nil},
			http.MethodGet,
			"/user/",
			http.StatusOK,
			responseAllUsers{Users: rs.users},
		},
		"404 - wrong verb": {
			[]any{},
			http.MethodPost,
			"/user/",
			http.StatusNotFound,
			responseNotFound{Message: "Page not found"},
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			rs.userMock.
				On("List").
				Return(tc.mockReturnArgs...).
				Once()

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(tc.method, tc.path, nil)
			rs.router.ServeHTTP(w, req)

			expectedBody, _ := json.Marshal(tc.expectedBody)

			assert.Equal(t, w.Code, tc.expectedStatus, "Wrong code received")
			assert.Equal(t, string(expectedBody), w.Body.String(), "Wrong response body")
		})
	}
}

func (rs *routerSuit) TestUserFetch() {
	t := rs.T()

	testCases := map[string]struct {
		mockInputArgs  []any
		mockReturnArgs []any
		method         string
		path           string
		expectedStatus int
		expectedBody   any
	}{
		"200 - Good call": {
			[]any{int(rs.users[1].ID)},
			[]any{rs.users[1], nil},
			http.MethodGet,
			fmt.Sprintf("/user/%d", int(rs.users[1].ID)),
			http.StatusOK,
			responseOneUser{User: rs.users[1]},
		},
		"200 - Good call, no user": {
			[]any{len(rs.users) + 1},
			[]any{entity.User{}, nil},
			http.MethodGet,
			fmt.Sprintf("/user/%d", len(rs.users)+1),
			http.StatusOK,
			responseOneUser{User: entity.User{}},
		},
		"400 - not a valid ID": {
			[]any{},
			[]any{},
			http.MethodGet,
			"/user/id",
			http.StatusBadRequest,
			responseError{Error: "Not a valid ID"},
		},
		"404 - wrong verb": {
			[]any{},
			[]any{},
			http.MethodPost,
			"/user/1",
			http.StatusNotFound,
			responseNotFound{Message: "Page not found"},
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			rs.userMock.
				On("Fetch", tc.mockInputArgs...).
				Return(tc.mockReturnArgs...).
				Once()

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(tc.method, tc.path, nil)
			rs.router.ServeHTTP(w, req)

			expectedBody, _ := json.Marshal(tc.expectedBody)

			assert.Equal(t, w.Code, tc.expectedStatus, "Wrong code received")
			assert.Equal(t, string(expectedBody), w.Body.String(), "Wrong response body")
		})
	}
}
