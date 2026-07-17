// Package apperrors defines the structured error type usecases and
// repositories return for any failure that should reach the client as a
// coded error rather than a bare 500. It mirrors the fixed error-code table
// in docs/api.md, so handlers can map any error to the right HTTP status and
// envelope without a per-entity switch statement.
package apperrors

import "net/http"

// Code is one of the fixed error codes clients are documented to handle.
type Code string

const (
	CodeValidation   Code = "validation_error"
	CodeUnauthorized Code = "unauthorized"
	CodeForbidden    Code = "forbidden"
	CodeNotFound     Code = "not_found"
	CodeConflict     Code = "conflict"
	CodeInternal     Code = "internal_error"
)

var statusByCode = map[Code]int{
	CodeValidation:   http.StatusBadRequest,
	CodeUnauthorized: http.StatusUnauthorized,
	CodeForbidden:    http.StatusForbidden,
	CodeNotFound:     http.StatusNotFound,
	CodeConflict:     http.StatusConflict,
	CodeInternal:     http.StatusInternalServerError,
}

// StatusFor returns the HTTP status for code, defaulting to 500 for an
// unknown or zero-value Code.
func StatusFor(code Code) int {
	if status, ok := statusByCode[code]; ok {
		return status
	}
	return http.StatusInternalServerError
}

// Error is the error type usecase/domain code should return for any failure
// that should reach the client as a structured error. Field is set only
// for CodeValidation. Err, if set, is the wrapped cause for server-side
// logging only — it's never serialized to the client.
type Error struct {
	Code    Code
	Message string
	Field   string
	Err     error
}

func (e *Error) Error() string {
	return e.Message
}

func (e *Error) Unwrap() error {
	return e.Err
}

// NotFound constructs a 404 apperrors.Error.
func NotFound(message string) *Error {
	return &Error{Code: CodeNotFound, Message: message}
}

// Conflict constructs a 409 apperrors.Error.
func Conflict(message string) *Error {
	return &Error{Code: CodeConflict, Message: message}
}

// Forbidden constructs a 403 apperrors.Error.
func Forbidden(message string) *Error {
	return &Error{Code: CodeForbidden, Message: message}
}

// Unauthorized constructs a 401 apperrors.Error.
func Unauthorized(message string) *Error {
	return &Error{Code: CodeUnauthorized, Message: message}
}

// Validation constructs a 400 apperrors.Error for the named field.
func Validation(field, message string) *Error {
	return &Error{Code: CodeValidation, Field: field, Message: message}
}

// Internal constructs a 500 apperrors.Error wrapping err. The client-facing
// message is always the generic "internal server error" — err is retained
// only so the caller can log it with request-ID correlation.
func Internal(err error) *Error {
	return &Error{Code: CodeInternal, Message: "internal server error", Err: err}
}
