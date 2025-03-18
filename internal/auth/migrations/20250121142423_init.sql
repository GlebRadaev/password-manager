-- +goose Up
-- +goose StatementBegin

CREATE SCHEMA IF NOT EXISTS auth;

DROP TABLE IF EXISTS auth.otp_codes;
DROP TABLE IF EXISTS auth.users;

CREATE TABLE IF NOT EXISTS auth.users
(
    user_id      UUID PRIMARY KEY,
    username     TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    email        TEXT UNIQUE NOT NULL,
    created_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS auth.otp_codes
(
    user_id    UUID REFERENCES auth.users(user_id) ON DELETE CASCADE,
    otp_code   TEXT NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    PRIMARY KEY (user_id, otp_code)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS auth.otp_codes;
DROP TABLE IF EXISTS auth.users;
DROP SCHEMA IF EXISTS auth;

-- +goose StatementEnd