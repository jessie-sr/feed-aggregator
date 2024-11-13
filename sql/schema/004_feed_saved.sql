-- +goose Up

CREATE TABLE feed_saved (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE, -- foreign key referencing "users" table; cascades delete
    feed_id UUID NOT NULL REFERENCES feeds(id) ON DELETE CASCADE, -- feed saved by the user
    UNIQUE(user_id, feed_id) -- shouldn't be able to save the same feed multiple times
);

-- +goose Down
DROP TABLE feed_saved;