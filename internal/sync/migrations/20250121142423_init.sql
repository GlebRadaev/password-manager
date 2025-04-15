-- +goose Up
-- +goose StatementBegin

CREATE SCHEMA sync;

CREATE TABLE sync.conflicts (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    data_id UUID NOT NULL,
    client_data BYTEA NOT NULL,
    server_data BYTEA NOT NULL,
    resolved BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_user_id ON sync.conflicts (user_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE sync.conflicts;
DROP SCHEMA sync;

-- +goose StatementEnd