package handlers

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/jha-captech/user-microservice/internal/models"
	"github.com/stretchr/testify/assert"

	serviceMock "github.com/jha-captech/user-microservice/internal/handlers/mock"
)

func TestHandleDeleteUser(t *testing.T) {
	mockService := new(serviceMock.MockUserDeleter)
	logger := slog.Default()
	handler := HandleDeleteUser(logger, mockService)

	tests := map[string]struct {
		mockFetchCalled  bool
		mockFetchInput   []any
		mockFetchOutput  []any
		mockDeleteCalled bool
		mockDeleteInput  []any
		mockDeleteOutput []any
		urlParam         string
		expectedCode     int
		expectedBody     string
	}{
		"valid request, user deleted": {
			mockFetchCalled:  true,
			mockFetchInput:   []any{1},
			mockFetchOutput:  []any{models.User{}, nil},
			mockDeleteCalled: true,
			mockDeleteInput:  []any{1},
			mockDeleteOutput: []any{nil},
			urlParam:         "1",
			expectedCode:     http.StatusAccepted,
			expectedBody:     toJSONString(responseMsg{Message: "object successful deleted"}),
		},
		"invalid ID format": {
			mockFetchCalled:  false,
			mockFetchInput:   nil,
			mockFetchOutput:  nil,
			mockDeleteCalled: false,
			mockDeleteInput:  nil,
			mockDeleteOutput: nil,
			urlParam:         "abc",
			expectedCode:     http.StatusBadRequest,
			expectedBody:     toJSONString(responseErr{Error: "Not a valid ID"}),
		},
		"user does not exist": {
			mockFetchCalled:  true,
			mockFetchInput:   []any{1},
			mockFetchOutput:  []any{models.User{}, sql.ErrNoRows},
			mockDeleteCalled: false,
			mockDeleteInput:  nil,
			mockDeleteOutput: nil,
			urlParam:         "1",
			expectedCode:     http.StatusBadRequest,
			expectedBody:     toJSONString(responseErr{Error: "Object does not exist"}),
		},
		"error deleting user": {
			mockFetchCalled:  true,
			mockFetchInput:   []any{1},
			mockFetchOutput:  []any{models.User{}, nil},
			mockDeleteCalled: true,
			mockDeleteInput:  []any{1},
			mockDeleteOutput: []any{errors.New("deletion error")},
			urlParam:         "1",
			expectedCode:     http.StatusInternalServerError,
			expectedBody:     toJSONString(responseErr{Error: "Error deleting object."}),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if tc.mockFetchCalled {
				mockService.
					On("FetchUser", tc.mockFetchInput...).
					Return(tc.mockFetchOutput...).
					Once()
			}
			if tc.mockDeleteCalled {
				mockService.
					On("DeleteUser", tc.mockDeleteInput...).
					Return(tc.mockDeleteOutput...).
					Once()
			}

			req, err := http.NewRequest(http.MethodDelete, "/api/user/"+tc.urlParam, nil)
			assert.NoError(t, err)

			// Add chi URLParam
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("ID", tc.urlParam)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tc.expectedCode, rr.Code, "Wrong code received")
			assert.JSONEq(t, tc.expectedBody, rr.Body.String(), "Wrong response body")

			if tc.mockFetchCalled {
				mockService.AssertExpectations(t)
			} else {
				mockService.AssertNotCalled(t, "FetchUser")
			}

			if tc.mockDeleteCalled {
				mockService.AssertExpectations(t)
			} else {
				mockService.AssertNotCalled(t, "DeleteUser")
			}
		})
	}
}
