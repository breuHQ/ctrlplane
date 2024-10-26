// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: github_repos.sql

package entities

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const createGitHubRepo = `-- name: CreateGitHubRepo :one
INSERT INTO github_repos (repo_id, installation_id, github_id, name, full_name, url, is_active)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id, created_at, updated_at, repo_id, installation_id, github_id, name, full_name, url, is_active
`

type CreateGitHubRepoParams struct {
	RepoID         uuid.UUID   `json:"repo_id"`
	InstallationID uuid.UUID   `json:"installation_id"`
	GithubID       int64       `json:"github_id"`
	Name           string      `json:"name"`
	FullName       string      `json:"full_name"`
	Url            string      `json:"url"`
	IsActive       pgtype.Bool `json:"is_active"`
}

func (q *Queries) CreateGitHubRepo(ctx context.Context, arg CreateGitHubRepoParams) (GithubRepo, error) {
	row := q.db.QueryRow(ctx, createGitHubRepo,
		arg.RepoID,
		arg.InstallationID,
		arg.GithubID,
		arg.Name,
		arg.FullName,
		arg.Url,
		arg.IsActive,
	)
	var i GithubRepo
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.RepoID,
		&i.InstallationID,
		&i.GithubID,
		&i.Name,
		&i.FullName,
		&i.Url,
		&i.IsActive,
	)
	return i, err
}

const deleteGitHubRepo = `-- name: DeleteGitHubRepo :one
DELETE FROM github_repos
WHERE id = $1
RETURNING id
`

func (q *Queries) DeleteGitHubRepo(ctx context.Context, id uuid.UUID) (uuid.UUID, error) {
	row := q.db.QueryRow(ctx, deleteGitHubRepo, id)
	err := row.Scan(&id)
	return id, err
}

const getGitHubRepoByFullName = `-- name: GetGitHubRepoByFullName :one
SELECT id, created_at, updated_at, repo_id, installation_id, github_id, name, full_name, url, is_active
FROM github_repos
WHERE full_name = $1
`

func (q *Queries) GetGitHubRepoByFullName(ctx context.Context, fullName string) (GithubRepo, error) {
	row := q.db.QueryRow(ctx, getGitHubRepoByFullName, fullName)
	var i GithubRepo
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.RepoID,
		&i.InstallationID,
		&i.GithubID,
		&i.Name,
		&i.FullName,
		&i.Url,
		&i.IsActive,
	)
	return i, err
}

const getGitHubRepoByID = `-- name: GetGitHubRepoByID :one
SELECT id, created_at, updated_at, repo_id, installation_id, github_id, name, full_name, url, is_active
FROM github_repos
WHERE id = $1
`

func (q *Queries) GetGitHubRepoByID(ctx context.Context, id uuid.UUID) (GithubRepo, error) {
	row := q.db.QueryRow(ctx, getGitHubRepoByID, id)
	var i GithubRepo
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.RepoID,
		&i.InstallationID,
		&i.GithubID,
		&i.Name,
		&i.FullName,
		&i.Url,
		&i.IsActive,
	)
	return i, err
}

const getGitHubRepoByInstallationIDAndGitHubID = `-- name: GetGitHubRepoByInstallationIDAndGitHubID :one
SELECT id, created_at, updated_at, repo_id, installation_id, github_id, name, full_name, url, is_active
FROM github_repos
WHERE installation_id = $1 AND github_id = $2
`

type GetGitHubRepoByInstallationIDAndGitHubIDParams struct {
	InstallationID uuid.UUID `json:"installation_id"`
	GithubID       int64     `json:"github_id"`
}

func (q *Queries) GetGitHubRepoByInstallationIDAndGitHubID(ctx context.Context, arg GetGitHubRepoByInstallationIDAndGitHubIDParams) (GithubRepo, error) {
	row := q.db.QueryRow(ctx, getGitHubRepoByInstallationIDAndGitHubID, arg.InstallationID, arg.GithubID)
	var i GithubRepo
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.RepoID,
		&i.InstallationID,
		&i.GithubID,
		&i.Name,
		&i.FullName,
		&i.Url,
		&i.IsActive,
	)
	return i, err
}

const getGitHubRepoByName = `-- name: GetGitHubRepoByName :one
SELECT id, created_at, updated_at, repo_id, installation_id, github_id, name, full_name, url, is_active
FROM github_repos
WHERE name = $1
`

func (q *Queries) GetGitHubRepoByName(ctx context.Context, name string) (GithubRepo, error) {
	row := q.db.QueryRow(ctx, getGitHubRepoByName, name)
	var i GithubRepo
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.RepoID,
		&i.InstallationID,
		&i.GithubID,
		&i.Name,
		&i.FullName,
		&i.Url,
		&i.IsActive,
	)
	return i, err
}

const getGitHubRepoByRepoID = `-- name: GetGitHubRepoByRepoID :many
SELECT id, created_at, updated_at, repo_id, installation_id, github_id, name, full_name, url, is_active
FROM github_repos
WHERE repo_id = $1
`

func (q *Queries) GetGitHubRepoByRepoID(ctx context.Context, repoID uuid.UUID) ([]GithubRepo, error) {
	rows, err := q.db.Query(ctx, getGitHubRepoByRepoID, repoID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GithubRepo
	for rows.Next() {
		var i GithubRepo
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.RepoID,
			&i.InstallationID,
			&i.GithubID,
			&i.Name,
			&i.FullName,
			&i.Url,
			&i.IsActive,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateGitHubRepo = `-- name: UpdateGitHubRepo :one
UPDATE github_repos
SET repo_id = $2, 
    installation_id = $3, 
    github_id = $4, 
    name = $5, 
    full_name = $6, 
    url = $7, 
    is_active = $8
WHERE id = $1
RETURNING id, created_at, updated_at, repo_id, installation_id, github_id, name, full_name, url, is_active
`

type UpdateGitHubRepoParams struct {
	ID             uuid.UUID   `json:"id"`
	RepoID         uuid.UUID   `json:"repo_id"`
	InstallationID uuid.UUID   `json:"installation_id"`
	GithubID       int64       `json:"github_id"`
	Name           string      `json:"name"`
	FullName       string      `json:"full_name"`
	Url            string      `json:"url"`
	IsActive       pgtype.Bool `json:"is_active"`
}

func (q *Queries) UpdateGitHubRepo(ctx context.Context, arg UpdateGitHubRepoParams) (GithubRepo, error) {
	row := q.db.QueryRow(ctx, updateGitHubRepo,
		arg.ID,
		arg.RepoID,
		arg.InstallationID,
		arg.GithubID,
		arg.Name,
		arg.FullName,
		arg.Url,
		arg.IsActive,
	)
	var i GithubRepo
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.RepoID,
		&i.InstallationID,
		&i.GithubID,
		&i.Name,
		&i.FullName,
		&i.Url,
		&i.IsActive,
	)
	return i, err
}