-- +goose Up

CREATE
EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE habits
(
    id          UUID PRIMARY KEY      DEFAULT gen_random_uuid(),
    user_id     UUID         NOT NULL,
    name        VARCHAR(255) NOT NULL,
    description TEXT,
    image_url   VARCHAR(500),
    created_at  TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_habits_user_id ON habits (user_id);

-- +goose Down

DROP TABLE IF EXISTS habits;