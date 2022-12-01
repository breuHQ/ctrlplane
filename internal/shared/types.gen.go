// Package shared provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.12.4 DO NOT EDIT.
package shared

// APIErrorResponse defines the structure of an API error response
type APIErrorResponse struct {
	// Code defines the code, helpful for debugging
	Code *int `json:"code,omitempty"`

	// Message defines the error to display to the user
	Message *string `json:"message,omitempty"`
}
