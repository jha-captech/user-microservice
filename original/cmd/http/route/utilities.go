package route

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// encodeResponse encodes a struct of type T as a JSON response.
func encodeResponse[T any](w http.ResponseWriter, status int, data T) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}
}

// decodeToStruct decodes a request body as a struct of type T.
func decodeToStruct[T any](r *http.Request) (T, error) {
	var data T
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return data, fmt.Errorf("decode json: %w", err)
	}
	return data, nil
}
