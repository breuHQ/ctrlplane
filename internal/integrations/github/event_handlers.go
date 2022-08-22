package github

import (
	"context"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"go.breu.io/ctrlplane/internal/cmn"
	"go.uber.org/zap"
)

// handles GitHub installation event. Below is the mermaid workflow.
//
//	sequenceDiagram
//	  autonumber
//	  actor UR as User
//	  participant UI as Browser
//	  participant GH as Github APP
//	  participant WH as API :: Webhook RX
//	  participant CI as API :: Comlete Installation
//	  participant WF as Workflow Engine
//	  participant DB
//	  UR ->> UI: Integrate Github
//	  UI ->> GH: Redirect to Github App Permissions Screen
//	  activate GH
//	    GH ->> WH: Receive Installation Data
//	      WH ->> WF: Send Installation Data to WF
//	      activate WF
//	    GH ->> UI: Receive Installation ID
//	  deactivate GH
//	  UI ->> CI: Send Installation ID
//	  activate CI
//	    CI ->> CI: Parse Team ID from Session
//	    CI ->> WF: Send to OnInstall workflow
//	    deactivate WF
//	  deactivate CI
//	  WF ->> DB: Save Installation
//	  activate WF
//	    WF ->> UI: Complete Installation
//	  deactivate WF
func handleInstallationEvent(ctx echo.Context) error {
	payload := &InstallationEventPayload{}
	if err := ctx.Bind(payload); err != nil {
		return err
	}

	workflows := &Workflows{}
	opts := cmn.Temporal.
		Queues[cmn.GithubIntegrationQueue].
		GetWorkflowOptions(strconv.Itoa(int(payload.Installation.ID)), string(InstallationEvent))

	run, err := cmn.Temporal.Client.
		SignalWithStartWorkflow(
			ctx.Request().Context(),
			opts.ID,
			InstallationEventSignal.String(),
			payload,
			opts,
			workflows.OnInstall,
			payload,
		)
	if err != nil {
		return nil
	}

	return ctx.JSON(http.StatusCreated, run.GetRunID())
}

// handles GitHub push event
func handlePushEvent(ctx echo.Context) error {
	payload := PushEventPayload{}
	if err := ctx.Bind(&payload); err != nil {
		cmn.Log.Info("Error: ", zap.Any("body", ctx.Request().Body), zap.String("error", err.Error()))
		return err
	}

	w := &Workflows{}
	opts := cmn.Temporal.
		Queues[cmn.GithubIntegrationQueue].
		GetWorkflowOptions(strconv.Itoa(int(payload.Installation.ID)), string(PushEvent), "ref", payload.After)

	exe, err := cmn.Temporal.Client.ExecuteWorkflow(context.Background(), opts, w.OnPush, payload)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusCreated, exe.GetRunID())
}
