-- +goose Up
-- +goose StatementBegin

-- Создание схемы data, если она не существует
CREATE SCHEMA IF NOT EXISTS data;

-- Создание таблицы data.data
CREATE TABLE IF NOT EXISTS data.data (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    type INT NOT NULL,
    data BYTEA NOT NULL,
    metadata JSONB, -- Используем JSONB вместо HSTORE
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Удаление таблицы data.data
DROP TABLE IF EXISTS data.data;

-- Удаление схемы data
DROP SCHEMA IF EXISTS data;

-- +goose StatementEnd