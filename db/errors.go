package db

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

// NewDatabase errors.
var (
	// ErrConnStringNotFound is returned when the connection string is not set.
	ErrConnStringNotSet = errors.New("connection string is not set")
)

var (
	ErrNewNoteCreateFailure = errors.New("failed to create new note")
)

var (
	// ErrInvalidInput is returned when the input is invalid.
	ErrInvalidInput = errors.New("invalid input")
	// ErrNotFound is returned when the resource is not found.
	ErrNotFound = errors.New("not found")
	// ErrAlreadyExists is returned when the resource already exists.
	ErrAlreadyExists = errors.New("already exists")
	// ErrInvalidID is returned when the ID is invalid.
	ErrInvalidID = errors.New("invalid ID")
)

// checkError checks and returns the appropriate error.
func checkError(err error) error {
	if err != nil {
		var responseError *azcore.ResponseError
		if errors.As(err, &responseError) {
			switch responseError.StatusCode {
			case http.StatusBadRequest:
				return ErrInvalidInput
			case http.StatusNotFound:
				return ErrNotFound
			case http.StatusConflict:
				return ErrAlreadyExists
			default:
				return ErrNewNoteCreateFailure
			}
		} else {
			return fmt.Errorf("%w: %w", ErrNewNoteCreateFailure, err)
		}
	}
	return ErrNewNoteCreateFailure
}
