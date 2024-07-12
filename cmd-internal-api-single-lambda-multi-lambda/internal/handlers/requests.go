package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// decodeRequestBody decodes a request body as a struct of type T.
func decodeRequestBody[T any](r *http.Request) (T, error) {
	var data T
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return data, fmt.Errorf("decode json: %w", err)
	}
	return data, nil
}
