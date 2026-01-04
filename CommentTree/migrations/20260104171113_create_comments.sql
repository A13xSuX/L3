-- +goose Up
CREATE TABLE IF NOT EXISTS comments(
    id bigserial PRIMARY KEY,
    parent_id bigint NULL REFERENCES comments(id) ON DELETE CASCADE,
    text varchar NOT NULL,
    created_at timestamptz NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_comments_parent_id ON comments(parent_id);
-- +goose Down
DROP TABLE IF EXISTS comments;