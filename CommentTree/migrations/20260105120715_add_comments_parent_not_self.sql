-- +goose Up
ALTER TABLE comments
    ADD CONSTRAINT comments_parent_not_self
        CHECK (parent_id IS NULL OR parent_id <> id);

-- +goose Down
ALTER TABLE comments
DROP CONSTRAINT IF EXISTS comments_parent_not_self;
