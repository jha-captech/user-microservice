package handler

import (
	"encoding/json"
	"fmt"
)

// structToJSON marshals a struct of type T to a JSON encoded string.
func structToJSON[T any](data T) (string, error) {
	JSONData, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("in structToJSON: %w", err)
	}
	return string(JSONData), nil
}
