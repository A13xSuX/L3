-- +goose Up
CREATE TABLE IF NOT EXISTS items(
    id BIGSERIAL PRIMARY KEY ,
    title VARCHAR(255) NOT NULL,
    sku VARCHAR(100) NOT NULL UNIQUE,
    quantity INT NOT NULL CHECK ( quantity>=0 ),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE TABLE IF NOT EXISTS users(
    id BIGSERIAL PRIMARY KEY ,
    username VARCHAR(255) NOT NULL UNIQUE ,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(20) NOT NULL DEFAULT 'viewer'
        CHECK (role IN ('manager', 'viewer', 'admin') ),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE TABLE IF NOT EXISTS audit(
    id BIGSERIAL PRIMARY KEY ,
    item_id BIGINT NOT NULL,
    action VARCHAR(20) NOT NULL CHECK (action IN ('UPDATE', 'DELETE', 'INSERT')),
    old_data JSONB,
    new_data JSONB,
    changed_by_user_id BIGINT,
    changed_by_username VARCHAR(255) NOT NULL,
    changed_by_role VARCHAR(20) NOT NULL DEFAULT 'viewer'
        CHECK (changed_by_role IN ('manager', 'viewer', 'admin') ),
    changed_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
-- +goose Down
DROP TABLE IF EXISTS audit;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS items;
