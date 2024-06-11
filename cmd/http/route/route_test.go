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
	"github.com/stretchr/testify/suite"

	"user-microservice/internal/database/entity"
)

// mocks

type serviceMock struct {
	Users []entity.User
}

func newServiceMock(userCount int) serviceMock {
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
	return serviceMock{Users: mockUsers}
}

func (sm serviceMock) List() ([]entity.User, error) {
	return sm.Users, nil
}

func (sm serviceMock) Fetch(ID int) (entity.User, error) {
	for _, mockUser := range sm.Users {
		if mockUser.ID != uint(ID) {
			continue
		}
		return mockUser, nil
	}
	return entity.User{}, nil
}

// test setup

type routerSuit struct {
	suite.Suite
	router       *gin.Engine
	userDataMock []entity.User
}

func TestRouterSuit(t *testing.T) {
	suite.Run(t, new(routerSuit))
}

func (rs *routerSuit) SetupTest() {
	userServiceMock := newServiceMock(5)
	rs.userDataMock = userServiceMock.Users
	rs.router = gin.Default()
	handler := NewHandler(userServiceMock)
	SetUpRoutes(rs.router, handler)
}

func structToJSONString[T any](data T) string {
	dataString, _ := json.Marshal(data)
	return string(dataString)
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
		method         string
		path           string
		expectedStatus int
		expectedBody   interface{}
	}{
		"200 - Good call": {
			http.MethodGet,
			"/user/",
			http.StatusOK,
			responseAllUsers{Users: rs.userDataMock},
		},
		"404 - wrong verb": {
			http.MethodPost,
			"/user/",
			http.StatusNotFound,
			responseNotFound{Message: "Page not found"},
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
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
		method         string
		path           string
		expectedStatus int
		expectedBody   interface{}
	}{
		"200 - Good call": {
			http.MethodGet,
			fmt.Sprintf("/user/%d", int(rs.userDataMock[1].ID)),
			http.StatusOK,
			responseOneUser{User: rs.userDataMock[1]},
		},
		"200 - Good call, no user": {
			http.MethodGet,
			fmt.Sprintf("/user/%d", len(rs.userDataMock)+1),
			http.StatusOK,
			responseOneUser{User: entity.User{}},
		},
		"400 - not a valid ID": {
			http.MethodGet,
			"/user/id",
			http.StatusBadRequest,
			responseError{Error: "Not a valid ID"},
		},
		"404 - wrong verb": {
			http.MethodPost,
			"/user/1",
			http.StatusNotFound,
			responseNotFound{Message: "Page not found"},
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(tc.method, tc.path, nil)
			rs.router.ServeHTTP(w, req)

			expectedBody, _ := json.Marshal(tc.expectedBody)

			assert.Equal(t, w.Code, tc.expectedStatus, "Wrong code received")
			assert.Equal(t, string(expectedBody), w.Body.String(), "Wrong response body")
		})
	}
}
