-- name: CreatePost :one
INSERT INTO posts (id, created_at, updated_at, title, description, published_at, url, feed_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
ON CONFLICT (url) DO UPDATE -- updates updated_at if the post's url already exists
SET updated_at = EXCLUDED.updated_at -- EXCLUDED: the row we want to insert
RETURNING *;