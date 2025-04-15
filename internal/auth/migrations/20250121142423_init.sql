-- +goose Up
-- +goose StatementBegin

CREATE SCHEMA auth;

CREATE TABLE auth.users
(
    user_id      UUID PRIMARY KEY,
    username     TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    email        TEXT UNIQUE NOT NULL,
    created_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE auth.otp_codes
(
    user_id    UUID REFERENCES auth.users(user_id) ON DELETE CASCADE,
    otp_code   TEXT NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    device_id  TEXT NOT NULL,
    PRIMARY KEY (user_id, otp_code, device_id)
);

CREATE TABLE auth.sessions
(
    session_id  UUID PRIMARY KEY,
    user_id     UUID REFERENCES auth.users(user_id) ON DELETE CASCADE,
    device_info TEXT NOT NULL,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    expires_at  TIMESTAMP NOT NULL
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE auth.otp_codes;
DROP TABLE auth.sessions;
DROP TABLE auth.users;
DROP SCHEMA auth;

-- +goose StatementEnd