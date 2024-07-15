package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/jha-captech/user-microservice/internal/models"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"

	serviceMock "github.com/jha-captech/user-microservice/internal/handlers/mock"
)

func toJSONString(data any) string {
	JSONString, err := json.Marshal(data)
	if err != nil {
		panic(fmt.Sprintf("Failed to Marshal data to JSON. \ndata: %v\nerr: %v", data, err))
	}
	return string(JSONString)
}

func TestHandleFetchUser(t *testing.T) {
	mockService := new(serviceMock.MockfetchUserServicer)
	logger := slog.Default()
	handler := HandleFetchUser(logger, mockService)

	users := []models.User{
		{ID: 1, FirstName: "John", LastName: "Doe", Role: "Admin", UserID: 1001},
		{ID: 2, FirstName: "Jane", LastName: "Smith", Role: "User", UserID: 1002},
	}

	usersOut := make([]outputUser, len(users))
	for i := 0; i < len(users); i++ {
		userOut := mapOutput(users[i])
		usersOut[i] = userOut
	}

	tests := map[string]struct {
		mockCalled     bool
		mockInput      []any
		mockOutput     []any
		requestIDParam string
		expectedCode   int
		expectedBody   string
	}{
		"valid ID, user found": {
			mockCalled:     true,
			mockInput:      []any{int(users[0].ID)},
			mockOutput:     []any{users[0], nil},
			requestIDParam: strconv.Itoa(int(users[0].ID)),
			expectedCode:   http.StatusOK,
			expectedBody:   toJSONString(responseUser{User: usersOut[0]}),
		},
		"invalid ID": {
			mockCalled:     false,
			mockInput:      nil,
			mockOutput:     nil,
			requestIDParam: "abc",
			expectedCode:   http.StatusBadRequest,
			expectedBody:   toJSONString(responseErr{Error: "Not a valid ID"}),
		},
		"user not found": {
			mockCalled:     true,
			mockInput:      []any{3},
			mockOutput:     []any{models.User{}, sql.ErrNoRows},
			requestIDParam: "3",
			expectedCode:   http.StatusOK,
			expectedBody:   toJSONString(responseUser{}),
		},
		"internal server error": {
			mockCalled:     true,
			mockInput:      []any{3},
			mockOutput:     []any{models.User{}, errors.New("")},
			requestIDParam: "3",
			expectedCode:   http.StatusInternalServerError,
			expectedBody:   toJSONString(responseErr{Error: "Internal server error"}),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if tc.mockCalled {
				mockService.
					On("FetchUser", tc.mockInput...).
					Return(tc.mockOutput...).
					Once()
			}

			req, err := http.NewRequest(http.MethodGet, "/api/user/"+tc.requestIDParam, nil)
			assert.NoError(t, err)

			// Add chi URLParam
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("ID", tc.requestIDParam)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tc.expectedCode, rr.Code, "Wrong code received")
			assert.JSONEq(t, tc.expectedBody, rr.Body.String(), "Wrong response body")

			if tc.mockCalled {
				mockService.AssertExpectations(t)
			} else {
				mockService.AssertNotCalled(t, "FetchUser")
			}
		})
	}
}
