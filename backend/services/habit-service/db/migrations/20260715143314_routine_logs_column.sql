-- +goose Up

CREATE TABLE routine_logs
(
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    routine_id   UUID      NOT NULL REFERENCES routines (id) ON DELETE CASCADE,
    completed_at TIMESTAMP NOT NULL,

    UNIQUE (routine_id, completed_at)
);

CREATE INDEX idx_routine_logs_routine_id
    ON routine_logs (routine_id);

CREATE INDEX idx_routine_logs_completed_at
    ON routine_logs (completed_at);

-- +goose Down

DROP TABLE IF EXISTS routine_logs;