package route

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"user-microservice/internal/database/entity"
	"user-microservice/internal/testutil"
)

// MOCKS

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

// TEST SETUP

type routerSuit struct {
	suite.Suite
	router   *chi.Mux
	userMock *serviceMock
	users    []entity.User
}

func TestRouterSuit(t *testing.T) {
	suite.Run(t, new(routerSuit))
}

func (rs *routerSuit) SetupSuite() {
	userServiceMock := new(serviceMock)
	rs.userMock = userServiceMock
	handler := NewHandler(userServiceMock)

	rs.router = chi.NewRouter()
	SetUpRoutes(rs.router, handler)
}

// TESTS

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

	users := testutil.NewUsers(3, testutil.WithIDStartRange(1))

	testCases := map[string]struct {
		mockReturnArgs []any
		method         string
		path           string
		expectedStatus int
		expectedBody   any
	}{
		"200 - Good call": {
			[]any{users, nil},
			http.MethodGet,
			"/user/",
			http.StatusOK,
			responseAllUsers{Users: users},
		},
		"404 - wrong verb": {
			[]any{},
			http.MethodPost,
			"/user/",
			http.StatusNotFound,
			responseMessage{Message: "Page not found"},
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

	user := testutil.NewUser()

	testCases := map[string]struct {
		mockInputArgs  []any
		mockReturnArgs []any
		method         string
		path           string
		expectedStatus int
		expectedBody   any
	}{
		"200 - Good call": {
			[]any{int(user.ID)},
			[]any{user, nil},
			http.MethodGet,
			fmt.Sprintf("/user/%d", int(user.ID)),
			http.StatusOK,
			responseOneUser{User: rs.users[1]},
		},
		"200 - Good call, no user": {
			[]any{int(user.ID) + 1},
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
			responseMessage{Message: "Page not found"},
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
