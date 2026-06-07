-- +goose Up

CREATE
EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE routines
(
    id         UUID PRIMARY KEY      DEFAULT gen_random_uuid(),
    user_id    UUID         NOT NULL,
    name       VARCHAR(255) NOT NULL,
    created_at TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_routines_user_id ON routines (user_id);

-- +goose Down

DROP TABLE IF EXISTS routines;