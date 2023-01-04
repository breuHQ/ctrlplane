// Package auth provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.12.4 DO NOT EDIT.
package auth

import (
	openapi_types "github.com/deepmap/oapi-codegen/pkg/types"
	"github.com/labstack/echo/v4"
	externalRef0 "go.breu.io/ctrlplane/internal/entity"
)

const (
	APIKeyAuthScopes = "APIKeyAuth.Scopes"
	BearerAuthScopes = "BearerAuth.Scopes"
)

// CreateAPIKeyRequest defines model for CreateAPIKeyRequest.
type CreateAPIKeyRequest struct {
	Name *string `json:"name,omitempty"`
}

// CreateAPIKeyResponse defines model for CreateAPIKeyResponse.
type CreateAPIKeyResponse struct {
	Key *string `json:"key,omitempty"`
}

// LoginRequest defines model for LoginRequest.
type LoginRequest struct {
	Email    openapi_types.Email `json:"email"`
	Password string              `json:"password"`
}

// RegisterationRequest defines model for RegisterationRequest.
type RegisterationRequest struct {
	ConfirmPassword string              `json:"confirm_password"`
	Email           openapi_types.Email `json:"email"`
	FirstName       string              `json:"first_name"`
	LastName        string              `json:"last_name"`
	Password        string              `json:"password"`
	TeamName        string              `json:"team_name"`
}

// RegisterationResponse defines model for RegisterationResponse.
type RegisterationResponse struct {
	Team *externalRef0.Team `json:"team,omitempty"`
	User *externalRef0.User `json:"user,omitempty"`
}

// TokenResponse defines model for TokenResponse.
type TokenResponse struct {
	AccessToken  *string `json:"access_token,omitempty"`
	RefreshToken *string `json:"refresh_token,omitempty"`
}

// ValidateAPIKeyResponse defines model for ValidateAPIKeyResponse.
type ValidateAPIKeyResponse struct {
	Message *string `json:"message,omitempty"`
}

// CreateTeamAPIKeyJSONRequestBody defines body for CreateTeamAPIKey for application/json ContentType.
type CreateTeamAPIKeyJSONRequestBody = CreateAPIKeyRequest

// CreateUserAPIKeyJSONRequestBody defines body for CreateUserAPIKey for application/json ContentType.
type CreateUserAPIKeyJSONRequestBody = CreateAPIKeyRequest

// LoginJSONRequestBody defines body for Login for application/json ContentType.
type LoginJSONRequestBody = LoginRequest

// RegisterJSONRequestBody defines body for Register for application/json ContentType.
type RegisterJSONRequestBody = RegisterationRequest

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Create a new API key for the team
	// (POST /auth/api-keys/team)
	CreateTeamAPIKey(ctx echo.Context) error

	// Create a new API key for the user
	// (POST /auth/api-keys/user)
	CreateUserAPIKey(ctx echo.Context) error

	// Validate an API key
	// (GET /auth/api-keys/validate)
	ValidateAPIKey(ctx echo.Context) error

	// Login
	// (POST /auth/login)
	Login(ctx echo.Context) error

	// Register a new user
	// (POST /auth/register)
	Register(ctx echo.Context) error

	// SecurityHandler returns the underlying Security Wrapper
	SecureHandler(handler echo.HandlerFunc, ctx echo.Context) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// CreateTeamAPIKey converts echo context to params.

func (w *ServerInterfaceWrapper) CreateTeamAPIKey(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{""})

	// Get the handler, get the secure handler if needed and then invoke with unmarshalled params.
	handler := w.Handler.CreateTeamAPIKey
	secure := w.Handler.SecureHandler
	err = secure(handler, ctx)

	return err
}

// CreateUserAPIKey converts echo context to params.

func (w *ServerInterfaceWrapper) CreateUserAPIKey(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{""})

	// Get the handler, get the secure handler if needed and then invoke with unmarshalled params.
	handler := w.Handler.CreateUserAPIKey
	secure := w.Handler.SecureHandler
	err = secure(handler, ctx)

	return err
}

// ValidateAPIKey converts echo context to params.

func (w *ServerInterfaceWrapper) ValidateAPIKey(ctx echo.Context) error {
	var err error

	ctx.Set(APIKeyAuthScopes, []string{""})

	// Get the handler, get the secure handler if needed and then invoke with unmarshalled params.
	handler := w.Handler.ValidateAPIKey
	secure := w.Handler.SecureHandler
	err = secure(handler, ctx)

	return err
}

// Login converts echo context to params.

func (w *ServerInterfaceWrapper) Login(ctx echo.Context) error {
	var err error

	// Get the handler, get the secure handler if needed and then invoke with unmarshalled params.
	handler := w.Handler.Login
	err = handler(ctx)

	return err
}

// Register converts echo context to params.

func (w *ServerInterfaceWrapper) Register(ctx echo.Context) error {
	var err error

	// Get the handler, get the secure handler if needed and then invoke with unmarshalled params.
	handler := w.Handler.Register
	err = handler(ctx)

	return err
}

// EchoRouter is an interface that wraps the methods of echo.Echo & echo.Group to provide a common interface
// for registering routes.
type EchoRouter interface {
	CONNECT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	TRACE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

// RegisterHandlers adds each server route to the EchoRouter.
func RegisterHandlers(router EchoRouter, si ServerInterface) {
	RegisterHandlersWithBaseURL(router, si, "")
}

// Registers handlers, and prepends BaseURL to the paths, so that the paths
// can be served under a prefix.
func RegisterHandlersWithBaseURL(router EchoRouter, si ServerInterface, baseURL string) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.POST(baseURL+"/auth/api-keys/team", wrapper.CreateTeamAPIKey)
	router.POST(baseURL+"/auth/api-keys/user", wrapper.CreateUserAPIKey)
	router.GET(baseURL+"/auth/api-keys/validate", wrapper.ValidateAPIKey)
	router.POST(baseURL+"/auth/login", wrapper.Login)
	router.POST(baseURL+"/auth/register", wrapper.Register)

}