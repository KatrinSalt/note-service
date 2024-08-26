package notes

import (
	"errors"
	"fmt"

	"github.com/KatrinSalt/notes-service/db"
)

// NewService errors.
var (
	// ErrDatabaseRequired is returned when the database instance is not provided.
	ErrDbRequired = errors.New("database is not provided")
	// ErrLoggerEmpty is returned when the logger instance is not provided.
	ErrLoggerRequired = errors.New("logger is not provided")
)

var (
	// ErrService is returned when the service fails.
	ErrService = errors.New("service error")
	// ErrInvalidInput is returned when the input is invalid.
	ErrInvalidInput = errors.New("invalid input")
	// ErrNotFound is returned when the resource is not found.
	ErrNotFound = errors.New("not found")
	// ErrCategoryNotFound = errors.New("category not found")
	// // ErrCategoryNotFound is returned when the category is not found.
	// ErrIDNotFound = errors.New("id not found")
	// ErrAlreadyExists is returned when the resource already exists.
	ErrAlreadyExists = errors.New("already exists")
)

// checkError checks and returns the appropriate error.
func checkError(err error) error {
	if err != nil {
		if errors.Is(err, db.ErrInvalidInput) {
			return ErrInvalidInput
		}
		if errors.Is(err, db.ErrNotFound) {
			return ErrNotFound
		}
		if errors.Is(err, db.ErrAlreadyExists) {
			return ErrAlreadyExists
		}
		return fmt.Errorf("%w: %w", ErrService, err)
	}
	return fmt.Errorf("%w: %w", ErrService, err)
}
