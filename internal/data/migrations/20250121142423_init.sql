-- +goose Up
-- +goose StatementBegin

CREATE SCHEMA IF NOT EXISTS data;

CREATE TABLE IF NOT EXISTS data.entries (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    type INT NOT NULL,
    data BYTEA NOT NULL,
    metadata JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS data.data;
DROP SCHEMA IF EXISTS data;

-- +goose StatementEnd