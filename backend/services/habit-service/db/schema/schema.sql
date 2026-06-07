CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE routines
(
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    UUID NOT NULL,
    name       VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE habits
(
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID NOT NULL,
    name        VARCHAR(255) NOT NULL,
    description TEXT,
    image_url   VARCHAR(500),
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE routine_habits
(
    routine_id UUID NOT NULL,
    habit_id   UUID NOT NULL,

    PRIMARY KEY (routine_id, habit_id),

    CONSTRAINT fk_routine_habits_routine
        FOREIGN KEY (routine_id)
            REFERENCES routines(id)
            ON DELETE CASCADE,

    CONSTRAINT fk_routine_habits_habit
        FOREIGN KEY (habit_id)
            REFERENCES habits(id)
            ON DELETE CASCADE
);

CREATE TABLE habit_logs
(
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    habit_id     UUID NOT NULL,
    completed_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_habit_logs_habit
        FOREIGN KEY (habit_id)
            REFERENCES habits(id)
            ON DELETE CASCADE,

    CONSTRAINT unique_habit_log
        UNIQUE (habit_id, completed_at)
);

CREATE INDEX idx_routines_user_id
    ON routines(user_id);

CREATE INDEX idx_habits_user_id
    ON habits(user_id);

CREATE INDEX idx_habit_logs_habit_id
    ON habit_logs(habit_id);