// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: user.sql

package entities

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const getUser = `-- name: GetUser :one

SELECT id, created_at, updated_at, org_id, email, first_name, last_name, password, is_active, is_verified
FROM users
WHERE id = $1
LIMIT 1
`

func (q *Queries) GetUser(ctx context.Context, id pgtype.UUID) (User, error) {
	row := q.db.QueryRow(ctx, getUser, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.OrgID,
		&i.Email,
		&i.FirstName,
		&i.LastName,
		&i.Password,
		&i.IsActive,
		&i.IsVerified,
	)
	return i, err
}