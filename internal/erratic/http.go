package erratic

import (
	"net/http"
)

// This are just alias to generate OpenAPI schema.
type (
	// BadRequestError represents a 400 Bad Request error.
	BadRequestError = QuantmError

	// UnauthorizedError represents a 401 Unauthorized error.
	UnauthorizedError = QuantmError

	ForbiddenError = QuantmError

	// NotFoundError represents a 404 Not Found error.
	NotFoundError = QuantmError
)

// NewBadRequestError creates a new 400 Bad Request error.
//
// Example:
//
//	err := NewBadRequestError("field", "invalid value")
//	fmt.Println(err.Error()) // Output: Bad Request
//	fmt.Println(err.Information) // Output: map[string]string{"field": "invalid value"}
func NewBadRequestError(args ...string) *QuantmError {
	return New(http.StatusBadRequest, "Bad Request.")
}

// NewUnauthorizedError creates a new 401 Unauthorized error.
//
// Example:
//
//	err := NewUnauthorizedError("user_id", "123")
//	fmt.Println(err.Error()) // Output: Unauthorized
//	fmt.Println(err.Information) // Output: map[string]string{"user_id": "123"}
func NewUnauthorizedError(args ...string) *QuantmError {
	return New(http.StatusUnauthorized, "Are you logged in?")
}

func NewForbiddenError(args ...string) *QuantmError {
	return New(http.StatusForbidden, "Permission Denied.")
}

// NewNotFoundError creates a new 404 Not Found error.
//
// Example:
//
//	err := NewNotFoundError("user_id", "123")
//	fmt.Println(err.Error()) // Output: Not Found
//	fmt.Println(err.Information) // Output: map[string]string{"user_id": "123"}
func NewNotFoundError(args ...string) *QuantmError {
	return New(http.StatusNotFound, "Resource Not Found.")
}

// NewInternalServerError creates a new 500 Internal Server Error.
//
// Example:
//
//	err := NewInternalServerError("reason", "database error")
//	fmt.Println(err.Error()) // Output: Internal Server Error
//	fmt.Println(err.Information) // Output: map[string]string{"reason": "database error"}
func NewInternalServerError(args ...string) *QuantmError {
	return New(500, "Internal Server Error.", args...)
}