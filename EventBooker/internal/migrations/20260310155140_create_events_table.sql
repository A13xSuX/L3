-- +goose Up
CREATE TABLE IF NOT EXISTS events (
    id UUID primary key DEFAULT gen_random_uuid(),
    title varchar(255) NOT NULL ,
    description TEXT DEFAULT '',
    date TIMESTAMP NOT NULL,
    total_seats int NOT NULL CHECK ( total_seats > 0 ),
    available_seats int NOT NULL CHECK (available_seats >= 0 and available_seats <= total_seats),
    price DECIMAL(10,2) NOT NULL DEFAULT 0.0,
    payment_required bool NOT NULL DEFAULT false,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS events;