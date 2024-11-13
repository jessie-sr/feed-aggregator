// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: feed_saved.sql

package db

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createFeedSaved = `-- name: CreateFeedSaved :one
INSERT INTO feed_saved (id, created_at, updated_at, user_id, feed_id)
VALUES ($1, $2, $3, $4, $5
)
RETURNING id, created_at, updated_at, user_id, feed_id
`

type CreateFeedSavedParams struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	UserID    uuid.UUID
	FeedID    uuid.UUID
}

func (q *Queries) CreateFeedSaved(ctx context.Context, arg CreateFeedSavedParams) (FeedSaved, error) {
	row := q.db.QueryRowContext(ctx, createFeedSaved,
		arg.ID,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.UserID,
		arg.FeedID,
	)
	var i FeedSaved
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.UserID,
		&i.FeedID,
	)
	return i, err
}

const getSavedFeeds = `-- name: GetSavedFeeds :many
SELECT id, created_at, updated_at, user_id, feed_id FROM feed_saved
`

func (q *Queries) GetSavedFeeds(ctx context.Context) ([]FeedSaved, error) {
	rows, err := q.db.QueryContext(ctx, getSavedFeeds)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []FeedSaved
	for rows.Next() {
		var i FeedSaved
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.UserID,
			&i.FeedID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
