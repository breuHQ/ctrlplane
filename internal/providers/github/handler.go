// Crafted with ❤ at Breu, Inc. <info@breu.io>, Copyright © 2022, 2024.
//
// Functional Source License, Version 1.1, Apache 2.0 Future License
//
// We hereby irrevocably grant you an additional license to use the Software under the Apache License, Version 2.0 that
// is effective on the second anniversary of the date we make the Software available. On or after that date, you may use
// the Software under the Apache License, Version 2.0, in which case the following will apply:
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with
// the License.
//
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on
// an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the
// specific language governing permissions and limitations under the License.

package github

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gocql/gocql"
	"github.com/labstack/echo/v4"

	"go.breu.io/quantm/internal/auth"
	"go.breu.io/quantm/internal/db"
	"go.breu.io/quantm/internal/shared"
)

type (
	ServerHandler struct{ *auth.SecurityHandler }
)

// NewServerHandler creates a new ServerHandler.
func NewServerHandler(middleware echo.MiddlewareFunc) *ServerHandler {
	return &ServerHandler{
		SecurityHandler: &auth.SecurityHandler{Middleware: middleware},
	}
}

func (s *ServerHandler) GithubWebhook(ctx echo.Context) error {
	signature := ctx.Request().Header.Get("X-Hub-Signature-256")

	if signature == "" {
		return ctx.JSON(http.StatusUnauthorized, ErrMissingHeaderGithubSignature)
	}

	// NOTE: We are reading the request body twice. This is not ideal.
	body, _ := io.ReadAll(ctx.Request().Body)
	ctx.Request().Body = io.NopCloser(bytes.NewBuffer(body))

	if err := Instance().VerifyWebhookSignature(body, signature); err != nil {
		return shared.NewAPIError(http.StatusUnauthorized, err)
	}

	headerEvent := ctx.Request().Header.Get("X-GitHub-Event")
	if headerEvent == "" {
		return shared.NewAPIError(http.StatusBadRequest, ErrMissingHeaderGithubEvent)
	}

	slog.Debug("GithubWebhook", "headerEvent", headerEvent)
	// Uncomment for debugging!
	// var jsonMap map[string]interface{}
	// json.Unmarshal([]byte(string(body)), &jsonMap)
	// slog.Debug("GithubWebhook", "body", jsonMap)

	event := WebhookEvent(headerEvent)
	handlers := WebhookEventHandlers{
		WebhookEventInstallation:             handleInstallationEvent,
		WebhookEventInstallationRepositories: handleInstallationRepositoriesEvent,
		WebhookEventPush:                     handlePushEvent,
		WebhookEventCreate:                   handleCreateOrDeleteEvent,
		WebhookEventDelete:                   handleCreateOrDeleteEvent,
		WebhookEventPullRequest:              handlePullRequestEvent,
		WebhookEventPullRequestReview:        handlePullRequestReviewEvent,
		WebhookEventPullRequestReviewComment: handlePullRequestReviewCommentEvent,
	}

	if handle, exists := handlers[event]; exists {
		return handle(ctx)
	} else {
		slog.Warn("Github Webhook: Unsupported event", "event", event)
	}

	return shared.NewAPIError(http.StatusBadRequest, ErrInvalidEvent)
}

func (s *ServerHandler) GithubCompleteInstallation(ctx echo.Context) error {
	request := &CompleteInstallationRequest{}
	if err := ctx.Bind(request); err != nil {
		return err
	}

	userID, _ := gocql.ParseUUID(ctx.Get("user_id").(string))
	payload := &CompleteInstallationSignal{request.InstallationID, request.SetupAction, userID}
	installation := &Installation{}
	workflows := &Workflows{}

	{
		opts := shared.Temporal().
			Queue(shared.ProvidersQueue).
			WorkflowOptions(
				shared.WithWorkflowBlock("github"),
				shared.WithWorkflowBlockID(strconv.Itoa(int(payload.InstallationID))),
				shared.WithWorkflowElement(WebhookEventInstallation.String()),
			)

		exe, err := shared.Temporal().
			Client().
			SignalWithStartWorkflow(
				ctx.Request().Context(),
				opts.ID,
				WorkflowSignalCompleteInstallation.String(),
				payload,
				opts,
				workflows.OnInstallationEvent,
			)
		if err != nil {
			return err
		}

		_ = exe.Get(ctx.Request().Context(), installation)
	}

	// TODO: handle this case!
	opts := shared.Temporal().Queue(shared.ProvidersQueue).
		WorkflowOptions(
			shared.WithWorkflowBlock("github"),
			shared.WithWorkflowBlockID(strconv.Itoa(int(payload.InstallationID))),
			shared.WithWorkflowElement(WebhookEventInstallation.String()),
			shared.WithWorkflowElementID("post-install"),
		)
	exe, err := shared.Temporal().Client().
		ExecuteWorkflow(
			ctx.Request().Context(),
			opts,
			workflows.PostInstall,
			installation,
		)

	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusCreated, &WorkflowResponse{RunID: exe.GetID(), Status: WorkflowStatusQueued})
}

func (s *ServerHandler) GithubGetInstallations(ctx echo.Context, params GithubGetInstallationsParams) error {
	result := make([]Installation, 0)
	filter := make(db.QueryParams)

	if params.InstallationId != nil {
		filter["installation_id"] = params.InstallationId.String()
	}

	if params.InstallationLogin != nil {
		filter["installation_login"] = shared.Quote(*params.InstallationLogin)
	}

	if params.InstallationId == nil && params.InstallationLogin == nil {
		filter["team_id"] = ctx.Get("team_id").(string)
	}

	if err := db.Filter(&Installation{}, &result, filter); err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, result)
}

func (s *ServerHandler) GithubGetRepos(ctx echo.Context) error {
	result := make([]Repo, 0)
	if err := db.Filter(
		&Repo{},
		&result,
		db.QueryParams{"team_id": ctx.Get("team_id").(string)},
	); err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, result)
}

func (s *ServerHandler) GithubListUserOrgs(ctx echo.Context, params GithubListUserOrgsParams) error {
	result := make([]OrgUser, 0)
	if err := db.Filter(&OrgUser{}, &result, db.QueryParams{"user_id": params.UserId}); err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, result)
}

func (s *ServerHandler) GithubCreateUserOrgs(ctx echo.Context) error {
	result := make([]OrgUser, 0)
	request := &CreateGithubUserOrgsRequest{}

	if err := ctx.Bind(request); err != nil {
		return err
	}

	if err := ctx.Validate(request); err != nil {
		return err
	}

	for _, id := range request.GithubOrgIDs {
		name := ""
		installation := &Installation{}

		err := db.Get(installation, db.QueryParams{"installation_id": id.String()})
		if err == nil {
			name = installation.InstallationLogin
		}

		orguser := &OrgUser{
			UserID:        request.UserID,
			GithubOrgID:   id,
			GithubUserID:  request.GithubUserID,
			GithubOrgName: name,
		}

		_ = db.Save(orguser) // TODO - update ORM to do a BulkSave Method.
		result = append(result, *orguser)
	}

	return ctx.JSON(http.StatusCreated, result)
}

func (s *ServerHandler) CreateTeamUser(ctx echo.Context) error {
	request := &CreateTeamUserRequest{}

	if err := ctx.Bind(request); err != nil {
		return shared.NewAPIError(http.StatusBadRequest, err)
	}

	if err := ctx.Validate(request); err != nil {
		return shared.NewAPIError(http.StatusBadRequest, err)
	}

	// associtaed user with team
	user := &auth.User{}
	if err := db.Get(user, db.QueryParams{"id": request.UserID.String()}); err != nil {
		slog.Error("get user", "debug", err.Error())
		return shared.NewAPIError(http.StatusNotFound, err)
	}

	user.TeamID = request.TeamID

	if err := db.Save(user); err != nil {
		slog.Error("save user", "debug", err.Error())
		return shared.NewAPIError(http.StatusInternalServerError, err)
	}

	orguser := &OrgUser{}
	filter := db.QueryParams{"github_org_id": request.GithubOrgID.String(), "github_user_id": request.GithubUserID.String()}

	if err := db.Get(orguser, filter); err != nil {
		slog.Error("get org user", "error", err.Error())
		return shared.NewAPIError(http.StatusNotFound, err)
	}

	orguser.UserID = request.UserID
	if err := db.Save(orguser); err != nil {
		slog.Error("update org user", "error", err.Error())
		return shared.NewAPIError(http.StatusInternalServerError, err)
	}

	// create teamuser
	teamuser := &auth.TeamUser{
		TeamID:                  request.TeamID,
		UserID:                  request.UserID,
		IsActive:                true,
		IsAdmin:                 false,
		IsMessageProviderLinked: false,
		UserLoginId:             orguser.GithubUserID,
	}
	if err := db.Save(teamuser); err != nil {
		slog.Error("create team user", "debug", err.Error())
		return shared.NewAPIError(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusCreated, user)
}
