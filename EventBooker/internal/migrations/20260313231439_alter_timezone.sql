-- +goose Up
ALTER TABLE bookings
    ALTER COLUMN created_at TYPE TIMESTAMPTZ
        USING created_at AT TIME ZONE 'Europe/Moscow',
    ALTER COLUMN expired_at TYPE TIMESTAMPTZ
        USING expired_at AT TIME ZONE 'Europe/Moscow',
    ALTER COLUMN confirmed_at TYPE TIMESTAMPTZ
        USING confirmed_at AT TIME ZONE 'Europe/Moscow';

-- +goose Down
ALTER TABLE bookings
    ALTER COLUMN created_at TYPE TIMESTAMP
        USING created_at AT TIME ZONE 'Europe/Moscow',
    ALTER COLUMN expired_at TYPE TIMESTAMP
        USING expired_at AT TIME ZONE 'Europe/Moscow',
    ALTER COLUMN confirmed_at TYPE TIMESTAMP
        USING confirmed_at AT TIME ZONE 'Europe/Moscow';