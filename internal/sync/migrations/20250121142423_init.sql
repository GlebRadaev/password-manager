-- +goose Up
-- +goose StatementBegin

CREATE SCHEMA IF NOT EXISTS sync;

CREATE TABLE IF NOT EXISTS sync.conflicts (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    data_id UUID NOT NULL,
    client_data BYTEA NOT NULL,
    server_data BYTEA NOT NULL,
    resolved BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_user_id ON sync.conflicts (user_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS sync.conflicts;
DROP SCHEMA IF EXISTS sync;

-- +goose StatementEnd