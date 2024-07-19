package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jha-captech/user-microservice/internal/models"
)

type inputUser struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Role      string `json:"role"`
	UserID    int    `json:"user_id"`
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

	// first name must not be blank
	if user.FirstName == "" {
		problems["first_name"] = "must not be blank"
	}

	// last name must not be blank
	if user.LastName == "" {
		problems["first_name"] = "must not be blank"
	}

	// validate role is `Customer` or `Employee`
	if user.Role == "" {
		problems["role"] = "must not be blank"
	} else if user.Role != "Customer" && user.Role != "Employee" {
		problems["role"] = "must be 'Customer' or 'Employee'"
	}

	// validate UserID greater than 0
	if user.UserID < 1 {
		problems["user_id"] = "must be more than 0"
	}

	return problems
}

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

func decodeValidateBody[I ValidatorMapper[O], O any](r *http.Request) (O, map[string]string, error) {
	var inputModel I

	// decode to JSON
	if err := json.NewDecoder(r.Body).Decode(&inputModel); err != nil {
		return *new(O), nil, fmt.Errorf("[in decodeValidateBody] decode json: %w", err)
	}

	// validate
	if problems := inputModel.Valid(); len(problems) > 0 {
		return *new(O), problems, fmt.Errorf(
			"[in decodeValidateBody] invalid %T: %d problems", inputModel, len(problems),
		)
	}

	// map to return type
	data, err := inputModel.MapTo()
	if err != nil {
		return *new(O), nil, fmt.Errorf("[in decodeValidateBody] map to %T: %w", *new(O), err)
	}

	return data, nil, nil
}
