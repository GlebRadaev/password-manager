-- +goose Up
-- +goose StatementBegin

CREATE SCHEMA IF NOT EXISTS sync;

CREATE TABLE IF NOT EXISTS sync.changes
(
    id         UUID PRIMARY KEY,
    user_id    UUID NOT NULL,
    data_id    UUID NOT NULL,
    type       TEXT NOT NULL,
    data       BYTEA NOT NULL,
    metadata   JSONB,
    timestamp  BIGINT NOT NULL
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS sync.changes;
DROP SCHEMA IF EXISTS sync;

-- +goose StatementEnd