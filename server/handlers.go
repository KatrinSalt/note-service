package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// encode writes the response as JSON to the response writer.
func encode[T any](w http.ResponseWriter, status int, v T) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		return fmt.Errorf("%w: %s", errors.New("malformed request body"), err)
	}
	return nil
}

// decode reads the request body as JSON and decodes it into the given value.
func decode[T any](r *http.Request) (T, error) {
	var v T
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		var syntaxError *json.SyntaxError
		if errors.As(err, &syntaxError) {
			err = ErrMalformedRequestBody
		}
		if errors.Is(err, io.EOF) {
			err = ErrEmptyRequestBody
		}
		return v, err
	}
	return v, nil
}

// logError creates a log message for error loggig.
func logError(err error, handler string) []any {
	return []any{"error", err, "type", "service", "handler", handler}
}
