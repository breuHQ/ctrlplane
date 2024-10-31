// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: messaging.sql

package entities

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
)

const createMessaging = `-- name: CreateMessaging :one
INSERT INTO messaging (hook, kind, link_to, data)
VALUES ($1, $2, $3, $4)
RETURNING id, created_at, updated_at, hook, kind, link_to, data
`

type CreateMessagingParams struct {
	Hook   string          `json:"hook"`
	Kind   string          `json:"kind"`
	LinkTo uuid.UUID       `json:"link_to"`
	Data   json.RawMessage `json:"data"`
}

func (q *Queries) CreateMessaging(ctx context.Context, arg CreateMessagingParams) (Messaging, error) {
	row := q.db.QueryRow(ctx, createMessaging,
		arg.Hook,
		arg.Kind,
		arg.LinkTo,
		arg.Data,
	)
	var i Messaging
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Hook,
		&i.Kind,
		&i.LinkTo,
		&i.Data,
	)
	return i, err
}

const getMessagesByLinkTo = `-- name: GetMessagesByLinkTo :many
SELECT id, created_at, updated_at, hook, kind, link_to, data
FROM messaging
WHERE link_to = $1
`

func (q *Queries) GetMessagesByLinkTo(ctx context.Context, linkTo uuid.UUID) ([]Messaging, error) {
	rows, err := q.db.Query(ctx, getMessagesByLinkTo, linkTo)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Messaging
	for rows.Next() {
		var i Messaging
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Hook,
			&i.Kind,
			&i.LinkTo,
			&i.Data,
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