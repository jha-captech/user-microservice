package handlers

import (
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/jha-captech/user-microservice/internal/models"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"

	serviceMock "github.com/jha-captech/user-microservice/internal/handlers/mock"
)

func TestHandleListUsers(t *testing.T) {
	mockService := new(serviceMock.MockUserLister)
	logger := slog.Default()
	handler := HandleListUsers(logger, mockService)

	users := []models.User{
		{ID: 1, FirstName: "John", LastName: "Doe", Role: "Admin", UserID: 1001},
		{ID: 2, FirstName: "Jane", LastName: "Smith", Role: "User", UserID: 1002},
	}

	usersOut := mapMultipleOutput(users)

	tests := map[string]struct {
		mockCalled   bool
		mockOutput   []any
		expectedCode int
		expectedBody string
	}{
		"users returned": {
			mockCalled:   true,
			mockOutput:   []any{users, nil},
			expectedCode: http.StatusOK,
			expectedBody: toJSONString(responseUsers{Users: usersOut}),
		},
		"no users found": {
			mockCalled:   true,
			mockOutput:   []any{[]models.User{}, nil},
			expectedCode: http.StatusOK,
			expectedBody: toJSONString(responseUsers{Users: []outputUser{}}),
		},
		"internal server error": {
			mockCalled:   true,
			mockOutput:   []any{[]models.User{}, errors.New("teat error")},
			expectedCode: http.StatusInternalServerError,
			expectedBody: toJSONString(responseErr{Error: "Error retrieving data"}),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if tc.mockCalled {
				mockService.
					On("ListUsers").
					Return(tc.mockOutput...).
					Once()
			}

			req, err := http.NewRequest(http.MethodGet, "/api/user", nil)
			assert.NoError(t, err)

			// Add chi URLParam
			rctx := chi.NewRouteContext()
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tc.expectedCode, rr.Code, "Wrong code received")
			assert.JSONEq(t, tc.expectedBody, rr.Body.String(), "Wrong response body")

			if tc.mockCalled {
				mockService.AssertExpectations(t)
			} else {
				mockService.AssertNotCalled(t, "ListUsers")
			}
		})
	}
}
