-- +goose Up
CREATE TABLE IF NOT EXISTS feeds (
    id UUID DEFAULT uuid_generate_v4(),
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    name TEXT NOT NULL,
    url TEXT NOT NULL UNIQUE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    PRIMARY KEY (id)
);

-- +goose Down
DROP TABLE feeds;
