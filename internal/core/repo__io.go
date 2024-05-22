// Copyright © 2023, Breu, Inc. <info@breu.io>. All rights reserved.
//
// This software is made available by Breu, Inc., under the terms of the BREU COMMUNITY LICENSE AGREEMENT, Version 1.0,
// found at https://www.breu.io/license/community. BY INSTALLING, DOWNLOADING, ACCESSING, USING OR DISTRIBUTING ANY OF
// THE SOFTWARE, YOU AGREE TO THE TERMS OF THE LICENSE AGREEMENT.
//
// The above copyright notice and the subsequent license agreement shall be included in all copies or substantial
// portions of the software.
//
// Breu, Inc. HEREBY DISCLAIMS ANY AND ALL WARRANTIES AND CONDITIONS, EXPRESS, IMPLIED, STATUTORY, OR OTHERWISE, AND
// SPECIFICALLY DISCLAIMS ANY WARRANTY OF MERCHANTABILITY OR FITNESS FOR A PARTICULAR PURPOSE, WITH RESPECT TO THE
// SOFTWARE.
//
// Breu, Inc. SHALL NOT BE LIABLE FOR ANY DAMAGES OF ANY KIND, INCLUDING BUT NOT LIMITED TO, LOST PROFITS OR ANY
// CONSEQUENTIAL, SPECIAL, INCIDENTAL, INDIRECT, OR DIRECT DAMAGES, HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY,
// ARISING OUT OF THIS AGREEMENT. THE FOREGOING SHALL APPLY TO THE EXTENT PERMITTED BY APPLICABLE LAW.

package core

import (
	"context"
	"time"

	"go.breu.io/quantm/internal/shared"
)

// RepoIO signals.
const (
	RepoIOSignalPush       shared.WorkflowSignal = "repo__push"
	ReopIOSignalRebase     shared.WorkflowSignal = "repo__push__rebase"
	RepoIOPullRequestLabel shared.WorkflowSignal = "repo__pull_request__label"
	RepoIOPullRequestMerge shared.WorkflowSignal = "repo__pull_request__merge"
)

// RepoIO signal payloads.
type (
	// RepoIO is the interface that defines the operations that can be performed on a repository.
	RepoIO interface {
		// GetRepoData gets the name & default branch for the provider repo.
		GetRepoData(ctx context.Context, id string) (*RepoIORepoData, error)

		// SetEarlyWarning sets the early warning flag for the provider repo.
		SetEarlyWarning(ctx context.Context, id string, value bool) error

		// GetAllBranches gets all the branches for the provider repo.
		GetAllBranches(ctx context.Context, payload *RepoIOInfoPayload) ([]string, error)

		DetectChanges(ctx context.Context, payload *RepoIODetectChangesPayload) (*RepoIOChanges, error)

		// TokenizedCloneURL returns the url with oauth token in it.
		//
		// NOTE - Since the url contains oauth token, it is best not to call this as activity.
		// LINK - https://github.com/orgs/community/discussions/24575#discussioncomment-3244524
		TokenizedCloneURL(ctx context.Context, payload *RepoIOInfoPayload) (string, error)

		// GetLatestCommit(ctx context.Context, payload *RepoIOGetLatestCommitPayload) (*LatestCommit, error)
		// DeployChangeset(ctx context.Context, payload *RepoIODeployChangesetPayload) error
		// TagCommit(ctx context.Context, payload *RepoIOTagCommitPayload) error
		// CreateBranch(ctx context.Context, payload *RepoIOCreateBranchPayload) error
		// DeleteBranch(ctx context.Context, payload *RepoIODeleteBranchPayload) error
		// MergeBranch(ctx context.Context, payload *RepoIOMergeBranchPayload) error
		// RebaseAndMerge(ctx context.Context, payload *RepoIORebaseAndMergePayload) (string, error)
		// DetectChange(ctx context.Context, payload *RepoIODetectChangePayload) (*BranchChanges, error)
		// GetAllBranches(ctx context.Context, payload *RepoIOGetAllBranchesPayload) ([]string, error)
		// TriggerCIAction(ctx context.Context, payload *RepoIOTriggerCIActionPayload) error
		// GetRepoTeamID(ctx context.Context, payload *RepoIOGetRepoTeamIDPayload) (string, error)
		// GetAllRelevantActions(ctx context.Context, payload *RepoIOGetAllRelevantActionsPayload) error
		// GetRepoByProviderID(ctx context.Context, payload *RepoIOGetRepoByProviderIDPayload) (*RepoProviderData, error)
		// UpdateRepoHasRarlyWarning(ctx context.Context, payload *RepoIOUpdateRepoHasRarlyWarningPayload) error
		// GetOrgUsers(ctx context.Context, payload *RepoIOGetOrgUsersPayload) error
		// RefreshDefaultBranches(ctx context.Context, payload *RepoIORefreshDefaultBranchesPayload) error
	}

	RepoSignalPushPayload struct {
		BranchRef      string         `json:"branch_ref"`
		Before         string         `json:"before"`
		After          string         `json:"after"`
		RepoName       string         `json:"repo_name"`
		RepoOwner      string         `json:"repo_owner"`
		CtrlID         string         `json:"ctrl_id"` // ID is the repo ID in the quantm DB. Should be UUID
		InstallationID shared.Int64   `json:"installation_id"`
		ProviderID     string         `json:"provider_id"`
		Commits        []RepoIOCommit `json:"commits"`
	}

	RepoSignalPullRequestLabelPayload struct{}

	RepoSignalPullRequestMergedPayload struct{}
)

// RepoIO types.
type (
	RepoIORepoData struct {
		Name          string `json:"name"`
		DefaultBranch string `json:"default_branch"`
		ProviderID    string `json:"provider_id"`
	}

	RepoIOInfoPayload struct {
		InstallationID shared.Int64 `json:"installation_id"`
		RepoName       string       `json:"repo_name"`
		RepoOwner      string       `json:"repo_owner"`
		DefaultBranch  string       `json:"defualt_branch"`
	}

	RepoIOClonePayload struct {
		Repo   *Repo                  `json:"repo"`   // Repo is the db record of the repo
		Push   *RepoSignalPushPayload `json:"push"`   // Push event payload
		Branch string                 `json:"branch"` // Branch to clone
		Path   string                 `json:"path"`   // Path to clone to
	}

	RepoIODetectChangesPayload struct {
		InstallationID shared.Int64 `json:"installation_id"`
		RepoName       string       `json:"repo_name"`
		RepoOwner      string       `json:"repo_owner"`
		DefaultBranch  string       `json:"defualt_branch"`
		TargetBranch   string       `json:"target_branch"`
	}

	RepoIOChanges struct {
		Added      shared.Int64 `json:"added"`
		Removed    shared.Int64 `json:"removed"`
		Modified   []string     `json:"modified"`
		Delta      shared.Int64 `json:"delta"`
		CompareUrl string       `json:"compare_url"`
		RepoUrl    string       `json:"repo_url"`
	}

	RepoIOCommit struct {
		SHA       string        `json:"sha"`
		Message   string        `json:"message"`
		Author    string        `json:"author"`
		Timestamp time.Time     `json:"timestamp"`
		Changes   RepoIOChanges `json:"changes"`
	}

	RepoIORebaseAtCommitResponse struct {
		SHA        string `json:"sha"`
		Message    string `json:"message"`
		InProgress bool   `json:"in_progress"`
	}
)