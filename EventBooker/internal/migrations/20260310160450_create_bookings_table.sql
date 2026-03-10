-- +goose Up
CREATE TABLE IF NOT EXISTS bookings(
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id uuid NOT NULL REFERENCES events(id) ON DELETE CASCADE ,
    username varchar(255) NOT NULL ,
    status varchar(30) NOT NULL CHECK (status IN('pending', 'confirmed', 'cancelled')),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    expired_at TIMESTAMP ,
    confirmed_at  TIMESTAMP
);
-- +goose Down
DROP TABLE IF EXISTS bookings;
