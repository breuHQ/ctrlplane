// Crafted with ❤ at Breu, Inc. <info@breu.io>, Copyright © 2024.
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

package code

import (
	"strings"
	"time"

	"github.com/gocql/gocql"
	"go.temporal.io/sdk/workflow"

	"go.breu.io/quantm/internal/core/defs"
	"go.breu.io/quantm/internal/db"
)

// BranchNameFromRef takes a full Git reference string and returns the branch name.
// For example, if the input is "refs/heads/my-branch", the output will be "my-branch".
func BranchNameFromRef(ref string) string {
	return strings.TrimPrefix(ref, "refs/heads/")
}

// RefFromBranchName takes a branch name and returns the full Git reference string.
// For example, if the input is "my-branch", the output will be "refs/heads/my-branch".
func RefFromBranchName(branch string) string {
	return "refs/heads/" + branch
}

// CreateQuantmRef takes a branch name and returns the full Git reference string for a quantum branch.
// For example, if the input is "my-branch", the output will be "refs/heads/quantm/my-branch".
func CreateQuantmRef(branch string) string {
	return "refs/heads/qtm/" + branch
}

// IsQuantmRef checks if a given Git reference string is a quantum branch reference.
// It returns true if the reference starts with "refs/heads/quantm/", otherwise false.
func IsQuantmRef(ref string) bool {
	return strings.HasPrefix(ref, "refs/heads/qtm/")
}

// IsQuantmBranch returns true if the given branch name starts with "qtm/".
// This is a helper function used to identify branches that are part of the Quantm project.
func IsQuantmBranch(branch string) bool {
	return strings.HasPrefix(branch, "qtm/")
}

// ToRebaseEvent handles the creation and setup of a rebase event.
func ToRebaseEvent(
	ctx workflow.Context,
	push *defs.Event[defs.Push, defs.RepoProvider],
	branch string,
	parent_id gocql.UUID,
) *defs.Event[defs.Rebase, defs.RepoProvider] {
	// Generate new UUID for the event
	id, _ := db.NewUUID()
	now := time.Now()

	// Create the rebase payload
	r := defs.Rebase{
		Ref:        push.Payload.Ref,
		Before:     push.Payload.Before,
		After:      push.Payload.After,
		HeadBranch: BranchNameFromRef(push.Payload.Ref),
		BaseBranch: branch,
		HeadCommit: *push.Payload.Commits.Latest(),
	}

	// Create a new rebase event
	rebase := &defs.Event[defs.Rebase, defs.RepoProvider]{
		Version: push.Version,
		ID:      id,
		Context: push.Context,
		Subject: push.Subject,
		Payload: r,
	}

	// Update the event with metadata
	rebase.SetParent(parent_id)
	rebase.SetScopeRebase()
	rebase.SetActionCreated()
	rebase.SetTimestamp(now)

	return rebase
}
