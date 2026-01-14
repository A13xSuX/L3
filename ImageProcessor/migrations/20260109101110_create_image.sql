-- +goose Up
CREATE TABLE IF NOT EXISTS images(
    id uuid PRIMARY KEY,
    status text NOT NULL,
    original_path text NOT NULL,
    processed_path text,
    thumb_path text,
    error text,
    created_at timestamptz NOT NULL DEFAULT NOW(),
    updated_at timestamptz NOT NULL DEFAULT NOW(),
    CONSTRAINT images_status_check CHECK (status IN ('queued','processing','ready','failed'))
);

-- +goose Down
DROP TABLE IF EXISTS images;