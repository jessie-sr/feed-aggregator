-- name: CreateFeedSaved :one
INSERT INTO feed_saved (id, created_at, updated_at, user_id, feed_id)
VALUES ($1, $2, $3, $4, $5
)
RETURNING *;

-- name: GetSavedFeeds :many
SELECT * FROM feed_saved; -- return all the feeds saved