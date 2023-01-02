// Package github provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.12.4 DO NOT EDIT.
package github

import (
	"encoding/json"
	"errors"

	"github.com/labstack/echo/v4"
)

const (
	APIKeyAuthScopes = "APIKeyAuth.Scopes"
	BearerAuthScopes = "BearerAuth.Scopes"
)

var (
	ErrInvalidSetupAction    = errors.New("invalid SetupAction value")
	ErrInvalidWorkflowStatus = errors.New("invalid WorkflowStatus value")
)

type (
	SetupActionMapType map[string]SetupAction // SetupActionMapType is a quick lookup map for SetupAction.
)

// Defines values for SetupAction.
const (
	SetupActionCreated SetupAction = "created"
	SetupActionDeleted SetupAction = "deleted"
	SetupActionUpdated SetupAction = "updated"
)

// SetupActionValues returns all known values for SetupAction.
var (
	SetupActionMap = SetupActionMapType{
		SetupActionCreated.String(): SetupActionCreated,
		SetupActionDeleted.String(): SetupActionDeleted,
		SetupActionUpdated.String(): SetupActionUpdated,
	}
)

/*
 * Helper methods for SetupAction for easy marshalling and unmarshalling.
 */
func (v SetupAction) String() string               { return string(v) }
func (v SetupAction) MarshalJSON() ([]byte, error) { return json.Marshal(v.String()) }
func (v *SetupAction) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	val, ok := SetupActionMap[s]
	if !ok {
		return ErrInvalidSetupAction
	}

	*v = val

	return nil
}

type (
	WorkflowStatusMapType map[string]WorkflowStatus // WorkflowStatusMapType is a quick lookup map for WorkflowStatus.
)

// Defines values for WorkflowStatus.
const (
	WorkflowStatusFailure  WorkflowStatus = "failure"
	WorkflowStatusQueued   WorkflowStatus = "queued"
	WorkflowStatusSignaled WorkflowStatus = "signaled"
	WorkflowStatusSkipped  WorkflowStatus = "skipped"
	WorkflowStatusSuccess  WorkflowStatus = "success"
)

// WorkflowStatusValues returns all known values for WorkflowStatus.
var (
	WorkflowStatusMap = WorkflowStatusMapType{
		WorkflowStatusFailure.String():  WorkflowStatusFailure,
		WorkflowStatusQueued.String():   WorkflowStatusQueued,
		WorkflowStatusSignaled.String(): WorkflowStatusSignaled,
		WorkflowStatusSkipped.String():  WorkflowStatusSkipped,
		WorkflowStatusSuccess.String():  WorkflowStatusSuccess,
	}
)

/*
 * Helper methods for WorkflowStatus for easy marshalling and unmarshalling.
 */
func (v WorkflowStatus) String() string               { return string(v) }
func (v WorkflowStatus) MarshalJSON() ([]byte, error) { return json.Marshal(v.String()) }
func (v *WorkflowStatus) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	val, ok := WorkflowStatusMap[s]
	if !ok {
		return ErrInvalidWorkflowStatus
	}

	*v = val

	return nil
}

// CompleteInstallationRequest complete the installation given the installation_id & setup_action
type CompleteInstallationRequest struct {
	InstallationId int64       `json:"installation_id"`
	SetupAction    SetupAction `json:"setup_action"`
}

// SetupAction defines model for SetupAction.
type SetupAction string

// WorkflowResponse workflow status & run id
type WorkflowResponse struct {
	RunId string `json:"run_id"`

	// Status the workflow status
	Status WorkflowStatus `json:"status"`
}

// WorkflowStatus the workflow status
type WorkflowStatus string

// GithubCompleteInstallationJSONRequestBody defines body for GithubCompleteInstallation for application/json ContentType.
type GithubCompleteInstallationJSONRequestBody = CompleteInstallationRequest

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Complete GitHub App installation
	// (POST /providers/github/complete-installation)
	GithubCompleteInstallation(ctx echo.Context) error

	// Get GitHub repositories
	// (GET /providers/github/repos)
	GithubGetRepos(ctx echo.Context) error

	// Webhook reciever for github
	// (POST /providers/github/webhook)
	GithubWebhook(ctx echo.Context) error

	// SecurityHandler returns the underlying Security Wrapper
	SecureHandler(handler echo.HandlerFunc, ctx echo.Context) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// GithubCompleteInstallation converts echo context to params.

func (w *ServerInterfaceWrapper) GithubCompleteInstallation(ctx echo.Context) error {
	var err error
	ctx.Set(BearerAuthScopes, []string{""})

	ctx.Set(APIKeyAuthScopes, []string{""})

	// Get the handler, get the secure handler if needed and then invoke with unmarshalled params.
	handler := w.Handler.GithubCompleteInstallation
	secure := w.Handler.SecureHandler
	err = secure(handler, ctx)

	return err
}

// GithubGetRepos converts echo context to params.

func (w *ServerInterfaceWrapper) GithubGetRepos(ctx echo.Context) error {
	var err error
	ctx.Set(BearerAuthScopes, []string{""})

	ctx.Set(APIKeyAuthScopes, []string{""})

	// Get the handler, get the secure handler if needed and then invoke with unmarshalled params.
	handler := w.Handler.GithubGetRepos
	secure := w.Handler.SecureHandler
	err = secure(handler, ctx)

	return err
}

// GithubWebhook converts echo context to params.

func (w *ServerInterfaceWrapper) GithubWebhook(ctx echo.Context) error {
	var err error

	// Get the handler, get the secure handler if needed and then invoke with unmarshalled params.
	handler := w.Handler.GithubWebhook
	err = handler(ctx)

	return err
}

// This is a simple interface which specifies echo.Route addition functions which
// are present on both echo.Echo and echo.Group, since we want to allow using
// either of them for path registration
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

	router.POST(baseURL+"/providers/github/complete-installation", wrapper.GithubCompleteInstallation)
	router.GET(baseURL+"/providers/github/repos", wrapper.GithubGetRepos)
	router.POST(baseURL+"/providers/github/webhook", wrapper.GithubWebhook)

}
