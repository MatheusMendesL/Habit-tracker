-- +goose Up

CREATE
EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE habit_logs
(
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    habit_id     UUID      NOT NULL REFERENCES habits (id) ON DELETE CASCADE,
    completed_at TIMESTAMP NOT NULL,

    UNIQUE (habit_id, completed_at)
);

CREATE INDEX idx_habit_logs_habit_id
    ON habit_logs (habit_id);

CREATE INDEX idx_habit_logs_completed_at
    ON habit_logs (completed_at);

-- +goose Down

DROP TABLE IF EXISTS habit_logs;