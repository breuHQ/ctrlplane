package defs

import (
	"go.breu.io/quantm/internal/shared"
)

type (
	// NOTE - this base struct need for any type of message. getting from core repo.
	MessageIOPayload struct {
		WorkspaceID string `json:"workspace_id"`
		ChannelID   string `json:"channel_id"`
		BotToken    string `json:"bot_token"`
		RepoName    string `json:"repo_name"`
		BranchName  string `json:"branch_name"`
		Author      string `json:"author"`
		AuthorUrl   string `json:"author_url"`
		IsChannel   bool   `json:"is_channel"`
	}

	// TODO: need to refine.
	MessageIOLineExeededPayload struct {
		MessageIOPayload *MessageIOPayload `json:"message_io_payload"`
		Threshold        shared.Int64      `json:"threshold"`
		DetectChanges    *RepoIOChanges    `json:"detect_changes"`
	}

	// TODO: need to refine.
	MessageIOMergeConflictPayload struct {
		MessageIOPayload *MessageIOPayload `json:"message_io_payload"`
		CommitUrl        string            `json:"commit_url"`
		RepoUrl          string            `json:"repo_url"`
		SHA              string            `json:"sha"`
	}

	// TODO: need to refine.
	MessageIOStaleBranchPayload struct {
		MessageIOPayload *MessageIOPayload `json:"message_io_payload"`
		CommitUrl        string            `json:"commit_url"`
		RepoUrl          string            `json:"repo_url"`
	}
)
