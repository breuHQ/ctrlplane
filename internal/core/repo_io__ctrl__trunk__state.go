package core

import (
	"log/slog"
	"time"

	"go.temporal.io/sdk/workflow"

	"go.breu.io/quantm/internal/shared"
)

type (
	// TrunkState defines the state for RepoWorkflows.TrunkCtrl.
	//
	// NOTE: This state is local to the workflow and all the members are private. It cannot be passed to child workflows.
	TrunkState struct {
		activties *RepoActivities     // activities is the activities for the branch control state
		repo      *Repo               // repo is the db record of the repo
		info      *RepoIOProviderInfo // info is the provider info for the repo
		branches  []string            // branches is the list of branches for the repo except the default branch.
		mutex     workflow.Mutex      // mutex is the mutex for the repo control state
		active    bool                // active is whether the branch control state is currently active
		counter   int                 // counter is the number of activity executions
	}
)

// on_push is a channel handler that processes push events for the repository.
//
// It receives a RepoIOSignalPushPayload from the push signal channel and signals the branch control state
// for the corresponding branch with the RepoIOSignalRebase signal and the received payload.
func (state *TrunkState) on_push(ctx workflow.Context) shared.ChannelHandler {
	return func(rx workflow.ReceiveChannel, more bool) {
		push := &RepoIOSignalPushPayload{}
		rx.Receive(ctx, push)

		for _, branch := range state.branches {
			if branch == BranchNameFromRef(push.BranchRef) {
				// NOTE - when to main branch skip the rebase workflow.
				// state.branches has all the branch including main.
				continue
			}

			// run rebase on all the feature branches. (except main).
			state.signal_branch(ctx, branch, RepoIOSignalRebase, push)
		}
	}
}

// on_create_delete is a channel handler that processes create or delete events for the repository.
//
// It receives a RepoIOSignalCreateOrDeletePayload from the create/delete signal channel and signals the branch control state
// for the corresponding branch with the RepoIOSignalCreateOrDelete signal and the received payload.
//
// If the payload indicates a branch was created, the function adds the branch to the list of branches in the state.
// If the payload indicates a branch was deleted, the function removes the branch from the list of branches in the state.
func (state *TrunkState) on_create_delete(ctx workflow.Context) shared.ChannelHandler {
	return func(rx workflow.ReceiveChannel, more bool) {
		create_delete := &RepoIOSignalCreateOrDeletePayload{}
		rx.Receive(ctx, create_delete)

		if create_delete.ForBranch(ctx) {
			if create_delete.IsCreated {
				state.add_branch(ctx, create_delete.Ref)
			} else {
				state.remove_branch(ctx, create_delete.Ref)
			}
		}
	}
}

// set_done marks the RepoCtrlState as inactive, releasing the mutex lock.
// This function should be called when the branch control state is no longer needed,
// such as when the branch is being deleted or merged.
func (state *TrunkState) set_done(ctx workflow.Context) {
	_ = state.mutex.Lock(ctx)
	defer state.mutex.Unlock()

	state.active = false
}

// is_active returns whether the RepoCtrlState is currently active.
// When the state is active, it means the branch control state is in use and the mutex is locked.
func (state *TrunkState) is_active() bool {
	return state.active
}

// refresh_info refreshes the provider information for the repository.
//
// It executes an activity to fetch the provider information and updates the state with the retrieved data.
func (state *TrunkState) refresh_info(ctx workflow.Context) {
	opts := workflow.ActivityOptions{StartToCloseTimeout: 60 * time.Second}
	ctx = workflow.WithActivityOptions(ctx, opts)
	io := Instance().RepoIO(state.repo.Provider)
	info := &RepoIOProviderInfo{}

	_ = state.do(ctx, "refresh_provider_info", io.GetProviderInfo, state.repo.CtrlID.String(), info, nil)

	_ = state.mutex.Lock(ctx)
	defer state.mutex.Unlock()

	state.info = info
}

// refresh_branches refreshes the list of branches for the repository.
//
// It executes an activity to fetch all branches from the provider and updates the state with the retrieved data.
func (state *TrunkState) refresh_branches(ctx workflow.Context) {
	opts := workflow.ActivityOptions{StartToCloseTimeout: 60 * time.Second}
	ctx = workflow.WithActivityOptions(ctx, opts)
	branches := make([]string, 0)
	io := Instance().RepoIO(state.repo.Provider)

	if state.info == nil {
		state.refresh_info(ctx)
	}

	_ = state.do(ctx, "refresh_branches", io.GetAllBranches, state.info, &branches)

	_ = state.mutex.Lock(ctx)
	defer state.mutex.Unlock()

	// TODO - pop the default branch from the branches array.

	state.branches = branches
}

// add_branch adds a branch to the list of branches in the state.
//
// It acquires the mutex lock, appends the branch to the list of branches, and releases the lock.
func (state *TrunkState) add_branch(ctx workflow.Context, branch string) {
	_ = state.mutex.Lock(ctx)
	defer state.mutex.Unlock()

	if branch != "" || branch != state.repo.DefaultBranch {
		state.branches = append(state.branches, branch)
	}
}

// remove_branch removes a branch from the list of branches in the state.
//
// It acquires the mutex lock, iterates through the list of branches, removes the specified branch, and releases the lock.
func (state *TrunkState) remove_branch(ctx workflow.Context, branch string) {
	_ = state.mutex.Lock(ctx)
	defer state.mutex.Unlock()

	for i, b := range state.branches {
		if b == branch {
			state.branches = append(state.branches[:i], state.branches[i+1:]...)
			break
		}
	}
}

// signal_branch sends a signal to the branch control state for the specified branch.
//
// It executes an activity to signal the branch control state with the specified signal and payload.
func (state *TrunkState) signal_branch(ctx workflow.Context, branch string, signal shared.WorkflowSignal, payload any) {
	opts := workflow.ActivityOptions{StartToCloseTimeout: 60 * time.Second}
	ctx = workflow.WithActivityOptions(ctx, opts)

	next := &RepoIOSignalBranchCtrlPayload{state.repo, branch, signal, payload}

	_ = state.do(
		ctx, "signal_branch_ctrl", state.activties.SignalBranch, next, nil,
		slog.String("signal", signal.String()),
	)
}

// increment is a helper function that increments the steps counter in the RepoIOBranchCtrlState.
//
// It acquires the mutex lock, increments the counter, and releases the lock.
func (state *TrunkState) increment(ctx workflow.Context) {
	_ = state.mutex.Lock(ctx)
	defer state.mutex.Unlock()

	state.counter++
}

// terminate is called to mark the RepoIOBranchCtrlState as inactive and cancel any associated timers.
// This function should be called when the branch control state is no longer needed, such as branch is being deleted or merged.
func (state *TrunkState) terminate(ctx workflow.Context) {
	state.set_done(ctx)
}

// do is a helper function that executes an activity within the context of the RepoIOBranchCtrlState.
// It logs the start and success of the activity, and increments the steps counter in the state.
// If the activity fails, the function logs the error and returns it.
//
// NOTE: This assumes that workflow.Context has been updated with activity options.
func (state *TrunkState) do(ctx workflow.Context, action string, activity, payload, result any, keyvals ...any) error {
	return _do(ctx, state.repo, state.repo.DefaultBranch, "trunk_ctrl", action, activity, payload, result, keyvals...)
}

// NewTrunkState creates a new TrunkState with the specified repo and activities.
func NewTrunkState(ctx workflow.Context, repo *Repo) *TrunkState {
	return &TrunkState{
		activties: &RepoActivities{},
		repo:      repo,
		branches:  make([]string, 0),
		mutex:     workflow.NewMutex(ctx),
		active:    true,
		counter:   0,
	}
}
