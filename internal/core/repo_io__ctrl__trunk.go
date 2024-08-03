package core

import (
	"go.temporal.io/sdk/workflow"
)

// TrunkCtrl is the event loop to process events during the lifecycle of the main branch.
//
// It processes the following events:
//
//   - push
//   - create_delete
func TrunkCtrl(ctx workflow.Context, repo *Repo) error {
	state := NewTrunkState(ctx, repo)
	selector := workflow.NewSelector(ctx)

	// channels
	// push event
	push := workflow.GetSignalChannel(ctx, RepoIOSignalPush.String())
	selector.AddReceive(push, state.on_push(ctx))

	// create_delete
	create_delete := workflow.GetSignalChannel(ctx, RepoIOSignalCreateOrDelete.String())
	selector.AddReceive(create_delete, state.on_create_delete(ctx))

	// main event loop
	for state.is_active() {
		selector.Select(ctx)

		if state.needs_reset() {
			return state.as_new(ctx, "event history exceeded threshold", TrunkCtrl, repo)
		}
	}

	// graceful shutdown
	state.terminate(ctx)

	return nil
}