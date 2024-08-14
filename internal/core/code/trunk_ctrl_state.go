package code

import (
	"go.temporal.io/sdk/workflow"

	"go.breu.io/quantm/internal/core/defs"
	"go.breu.io/quantm/internal/shared"
)

type (
	TrunkState struct {
		*BaseCtrl
		active_branch string
	}
)

func (state *TrunkState) on_push(ctx workflow.Context) shared.ChannelHandler {
	return func(rx workflow.ReceiveChannel, more bool) {
		push := &defs.RepoIOSignalPushPayload{}
		state.rx(ctx, rx, push)

		for _, branch := range state.branches {
			if branch == BranchNameFromRef(push.BranchRef) {
				continue
			}

			state.signal_branch(ctx, branch, defs.RepoIOSignalRebase, push)
		}
	}
}

func (state *TrunkState) on_create_delete(ctx workflow.Context) shared.ChannelHandler {
	return func(rx workflow.ReceiveChannel, more bool) {
		create_delete := &defs.RepoIOSignalCreateOrDeletePayload{}
		state.rx(ctx, rx, create_delete)

		if create_delete.ForBranch(ctx) {
			if create_delete.IsCreated {
				state.add_branch(ctx, create_delete.Ref)
			} else {
				state.remove_branch(ctx, create_delete.Ref)
			}
		}
	}
}

func NewTrunkState(ctx workflow.Context, repo *defs.Repo) *TrunkState {
	return &TrunkState{
		BaseCtrl:      NewBaseCtrl(ctx, "trunk_ctrl", repo),
		active_branch: repo.DefaultBranch,
	}
}