package route

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

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

func (sm *serviceMock) Update(ID int, user entity.User) (entity.User, error) {
	args := sm.Called(ID, user)
	return args.Get(0).(entity.User), args.Error(1)
}

func (sm *serviceMock) Create(user entity.User) (int, error) {
	args := sm.Called(user)
	return args.Int(0), args.Error(1)
}

func (sm *serviceMock) Delete(ID int) error {
	args := sm.Called(ID)
	return args.Error(0)
}

// ━━ TEST SETUP ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

type routerSuit struct {
	suite.Suite
	router   *chi.Mux
	userMock *serviceMock
}

func TestRouterSuit(t *testing.T) {
	suite.Run(t, new(routerSuit))
}

func (rs *routerSuit) SetupSuite() {
	userServiceMock := new(serviceMock)
	rs.userMock = userServiceMock

	logger := slog.Default()

	handler := NewHandler(userServiceMock, logger)

	rs.router = chi.NewRouter()
	rs.router.Use(middleware.Logger)
	SetUpRoutes(rs.router, handler)
}

// ━━ TESTS ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

func (rs *routerSuit) TestHealthCheck() {
	t := rs.T()

	testCases := map[string]struct {
		method         string
		path           string
		expectedStatus int
	}{
		"200": {
			http.MethodGet,
			"/api/health-check",
			http.StatusOK,
		},
		"405": {
			http.MethodPut,
			"/api/health-check",
			http.StatusMethodNotAllowed,
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
			"/api/user/",
			http.StatusOK,
			responseAllUsers{Users: users},
		},
		"405 - wrong verb": {
			[]any{},
			http.MethodPost,
			"/api/user/",
			http.StatusMethodNotAllowed,
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

			assert.Equal(t, tc.expectedStatus, w.Code, "Wrong code received")
			assert.Equal(
				t,
				string(expectedBody),
				strings.TrimSpace(w.Body.String()),
				"Wrong response body",
			)
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
			fmt.Sprintf("/api/user/%d", int(user.ID)),
			http.StatusOK,
			responseOneUser{User: user},
		},
		"200 - Good call, no user": {
			[]any{int(user.ID) + 1},
			[]any{entity.User{}, nil},
			http.MethodGet,
			fmt.Sprintf("/api/user/%d", int(user.ID)+1),
			http.StatusOK,
			responseOneUser{User: entity.User{}},
		},
		"400 - not a valid ID": {
			[]any{},
			[]any{},
			http.MethodGet,
			"/api/user/id",
			http.StatusBadRequest,
			responseError{Error: "Not a valid ID"},
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

			assert.Equal(t, tc.expectedStatus, w.Code, "Wrong code received")
			assert.Equal(
				t,
				string(expectedBody),
				strings.TrimSpace(w.Body.String()),
				"Wrong response body",
			)
		})
	}
}

func (rs *routerSuit) TestUserUpdate() {
	t := rs.T()

	user := testutil.NewUser()

	testCases := map[string]struct {
		mockInputArgs  []any
		mockReturnArgs []any
		method         string
		path           string
		body           entity.User
		expectedStatus int
		expectedBody   any
	}{
		"200 - Good call": {
			[]any{int(user.ID), user},
			[]any{user, nil},
			http.MethodPut,
			fmt.Sprintf("/api/user/%d", int(user.ID)),
			user,
			http.StatusOK,
			responseOneUser{User: user},
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			rs.userMock.
				On("Update", tc.mockInputArgs...).
				Return(tc.mockReturnArgs...).
				Once()

			bodyJSON, _ := json.Marshal(tc.body)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(tc.method, tc.path, bytes.NewReader(bodyJSON))
			rs.router.ServeHTTP(w, req)

			expectedBody, _ := json.Marshal(tc.expectedBody)

			assert.Equal(t, tc.expectedStatus, w.Code, "Wrong code received")
			assert.Equal(
				t,
				string(expectedBody),
				strings.TrimSpace(w.Body.String()),
				"Wrong response body",
			)
		})
	}
}

func (rs *routerSuit) TestUserCreate() {
	t := rs.T()

	user := testutil.NewUser()

	testCases := map[string]struct {
		mockInputArgs  []any
		mockReturnArgs []any
		method         string
		path           string
		body           entity.User
		expectedStatus int
		expectedBody   any
	}{
		"200 - Good call": {
			[]any{user},
			[]any{int(user.ID), nil},
			http.MethodPost,
			fmt.Sprintf("/api/user"),
			user,
			http.StatusOK,
			responseID{ObjectID: int(user.ID)},
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			rs.userMock.
				On("Create", tc.mockInputArgs...).
				Return(tc.mockReturnArgs...).
				Once()

			bodyJSON, _ := json.Marshal(tc.body)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(tc.method, tc.path, bytes.NewReader(bodyJSON))
			rs.router.ServeHTTP(w, req)

			expectedBody, _ := json.Marshal(tc.expectedBody)

			assert.Equal(t, tc.expectedStatus, w.Code, "Wrong code received")
			assert.Equal(
				t,
				string(expectedBody),
				strings.TrimSpace(w.Body.String()),
				"Wrong response body",
			)
		})
	}
}

func (rs *routerSuit) TestDeleteCreate() {
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
			[]any{nil},
			http.MethodDelete,
			fmt.Sprintf("/api/user/%d", int(user.ID)),
			http.StatusOK,
			responseMessage{Message: "object successful deleted"},
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			rs.userMock.
				On("Delete", tc.mockInputArgs...).
				Return(tc.mockReturnArgs...).
				Once()

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(tc.method, tc.path, nil)
			rs.router.ServeHTTP(w, req)

			expectedBody, _ := json.Marshal(tc.expectedBody)

			assert.Equal(t, tc.expectedStatus, w.Code, "Wrong code received")
			assert.Equal(
				t,
				string(expectedBody),
				strings.TrimSpace(w.Body.String()),
				"Wrong response body",
			)
		})
	}
}
