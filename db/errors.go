package db

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

// NewDatabase errors.
var (
	// ErrConnStringNotFound is returned when the connection string is not provided.
	// ErrConnStringRequired = errors.New("connection string is not provided")
	// ErrDbIdEmpty is returned when the database id is not provided.
	// ErrDbIdRequired = errors.New("database id is not provided")
	// ErrContainerIdEmpty is returned when the container id is not provided.
	// ErrContainerIdRequired = errors.New("container id is not provided")
	// ErrLoggerEmpty is returned when the logger instance is not provided.
	ErrLoggerRequired = errors.New("logger is not provided")
	// ErrLoggerEmpty is returned when the logger instance is not provided.
	ErrClientRequired = errors.New("database client is not provided")
)

var (
	ErrClientConnection = errors.New("connection to the database failed")
)

// Generic error for the DB layer.
var (
	ErrInternalDB = errors.New("internal database error")
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
				return fmt.Errorf("%w: %w", ErrInternalDB, err)
			}
		} else {
			return fmt.Errorf("%w: %w", ErrInternalDB, err)
		}
	}
	return ErrInternalDB
}
