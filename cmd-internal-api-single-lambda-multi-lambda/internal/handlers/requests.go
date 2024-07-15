package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jha-captech/user-microservice/internal/models"
)

type Validator interface {
	Valid() (problems map[string]string)
}

type Mapper[T any] interface {
	MapTo() (T, error)
}

type ValidatorMapper[T any] interface {
	Validator
	Mapper[T]
}

type inputUser struct {
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Role      string `json:"role,omitempty"`
	UserID    int    `json:"user_id,omitempty"`
}

func (user inputUser) MapTo() (models.User, error) {
	return models.User{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      user.Role,
		UserID:    uint(user.UserID),
	}, nil
}

func (user inputUser) Valid() map[string]string {
	problems := make(map[string]string)

	// validate UserID greater than 0
	if user.UserID < 1 {
		problems["UserID"] = "UserID must be more than 0"
	}

	// validate role is `Customer` or `Employee`
	if user.Role != "Customer" && user.Role != "Employee" {
		problems["Role"] = "Role must be 'Customer' or 'Employee'"
	}

	return problems
}

func decodeValidateBody[I ValidatorMapper[O], O any](r *http.Request) (O, map[string]string, error) {
	var v I

	// decode to JSON
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return *new(O), nil, fmt.Errorf("decode json: %w", err)
	}

	// validate
	if problems := v.Valid(); len(problems) > 0 {
		return *new(O), problems, fmt.Errorf("invalid %I: %d problems", v, len(problems))
	}

	// map to return type
	data, err := v.MapTo()
	if err != nil {
	}

	return data, nil, nil
}
