package handlers

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/jha-captech/user-microservice/internal/models"
	"github.com/stretchr/testify/assert"

	serviceMock "github.com/jha-captech/user-microservice/internal/handlers/mock"
)

func TestHandleCreateUser(t *testing.T) {
	mockService := new(serviceMock.MockUserCreator)
	logger := slog.Default()
	handler := HandleCreateUser(logger, mockService)

	userIn := inputUser{FirstName: "John", LastName: "Doe", Role: "Customer", UserID: 1001}
	user := models.User{FirstName: "John", LastName: "Doe", Role: "Customer", UserID: 1001}

	tests := map[string]struct {
		mockCalled   bool
		mockInput    []any
		mockOutput   []any
		requestBody  string
		expectedCode int
		expectedBody string
	}{
		"valid request, user created": {
			mockCalled:   true,
			mockInput:    []any{user},
			mockOutput:   []any{1, nil},
			requestBody:  toJSONString(userIn),
			expectedCode: http.StatusCreated,
			expectedBody: toJSONString(responseID{ObjectID: 1}),
		},
		"invalid request body": {
			mockCalled:   false,
			mockInput:    nil,
			mockOutput:   nil,
			requestBody:  `{"FirstName":"John","LastName":"Doe","Role":"Admin"}`,
			expectedCode: http.StatusBadRequest,
			expectedBody: toJSONString(responseErr{
				ValidationErrors: map[string]string{
					"first_name": "must not be blank",
					"role":       "must be 'Customer' or 'Employee'",
					"user_id":    "must be more than 0",
				},
			}),
		},
		"error creating user": {
			mockCalled:   true,
			mockInput:    []any{user},
			mockOutput:   []any{0, errors.New("creation error")},
			requestBody:  toJSONString(userIn),
			expectedCode: http.StatusInternalServerError,
			expectedBody: toJSONString(responseErr{Error: "Error creating object"}),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodPost, "/api/user", strings.NewReader(tc.requestBody))
			assert.NoError(t, err)

			// Add chi URLParam
			rctx := chi.NewRouteContext()
			ctx := context.WithValue(req.Context(), chi.RouteCtxKey, rctx)
			req = req.WithContext(ctx)

			if tc.mockCalled {
				mockService.
					On("CreateUser", append([]any{ctx}, tc.mockInput...)...).
					Return(tc.mockOutput...).
					Once()
			}

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tc.expectedCode, rr.Code, "Wrong code received")
			assert.JSONEq(t, tc.expectedBody, rr.Body.String(), "Wrong response body")

			if tc.mockCalled {
				mockService.AssertExpectations(t)
			} else {
				mockService.AssertNotCalled(t, "CreateUser")
			}
		})
	}
}
