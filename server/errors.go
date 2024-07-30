package server

import (
	"errors"
	"net/http"

	"github.com/KatrinSalt/notes-service/notes"
)

var (
	// ErrInvalidRequest is returned when a request is invalid.
	ErrInvalidRequest = errors.New("invalid request")
	// ErrMalformedRequestBody is returned when the request body is malformed.
	ErrMalformedRequestBody = errors.New("malformed request body")
	// ErrEmptyRequestBody is returned when the request body is empty.
	ErrEmptyRequestBody = errors.New("empty request body")
	// ErrForbidden is returned when the request is forbidden.
	ErrForbidden = errors.New("forbidden")
	// ErrCategoryRequired is returned when a category is required.
	ErrCategoryRequired = errors.New("category is required")

	// // ErrIDRequired is returned when an id is required.
	// ErrIDRequired = errors.New("id is required")
)

// responseError is a response error.
type responseError struct {
	StatusCode int    `json:"statusCode"`
	Code       string `json:"code"`
	Message    string `json:"message"`
}

// Error returns the error message.
func (e *responseError) Error() string {
	return e.Message
}

// newResponseError creates a new response error with the given status code and error.
func newResponseError(statusCode int, code string, err error) error {
	return &responseError{
		StatusCode: statusCode,
		Code:       code,
		Message:    err.Error(),
	}
}

const (
	// CodeServerError is the error code for server error.
	CodeServerError = "ServerError"
	// CodeIDRequired is the error code for ID required.
	CodeIDRequired = "IDRequired"
	// CodeCategoryRequired is the error code for category required.
	CodeCategoryRequired = "CategoryRequired"
)

// errorCodeMaps contains a map with HTTP status codes and a map with errors
// and their codes.
var errorCodeMaps = map[int]map[error]string{
	http.StatusBadRequest: {
		ErrInvalidRequest:       "InvalidRequest",
		ErrMalformedRequestBody: "MalformedRequestBody",
		ErrEmptyRequestBody:     "EmptyRequestBody",
		notes.ErrInvalidInput:   "InvalidInput",
	},
	http.StatusNotFound: {
		notes.ErrNotFound: "NotFound",
	},
	http.StatusConflict: {
		notes.ErrAlreadyExists: "AlreadyExists",
	},
}

// errorCodes returns the status and error code for the given error.
func errorCodes(err error) (int, string) {
	for statusCode, errs := range errorCodeMaps {
		for e, code := range errs {
			if errors.Is(err, e) {
				return statusCode, code
			}
		}
	}
	return 0, ""
}

// writeError writes an error response to the caller.
func writeError(w http.ResponseWriter, statusCode int, code string, err error) {
	if err == nil {
		err = errors.New("internal server error")
	}

	respErr := newResponseError(statusCode, code, err)
	if err := encode(w, statusCode, respErr); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// writeServerError writes a server error response to the caller.
// Used to make sure the caller does not get any information about the
// internal error.
func writeServerError(w http.ResponseWriter) {
	writeError(w, http.StatusInternalServerError, CodeServerError, nil)
}
