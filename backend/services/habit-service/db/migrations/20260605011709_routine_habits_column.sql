-- +goose Up

CREATE TABLE routine_habits
(
    routine_id UUID NOT NULL,
    habit_id   UUID NOT NULL,

    PRIMARY KEY (routine_id, habit_id),

    CONSTRAINT fk_routine_habits_routine
        FOREIGN KEY (routine_id)
            REFERENCES routines (id)
            ON DELETE CASCADE,

    CONSTRAINT fk_routine_habits_habit
        FOREIGN KEY (habit_id)
            REFERENCES habits (id)
            ON DELETE CASCADE
);

CREATE INDEX idx_routine_habits_routine_id
    ON routine_habits (routine_id);

CREATE INDEX idx_routine_habits_habit_id
    ON routine_habits (habit_id);

-- +goose Down

DROP TABLE IF EXISTS routine_habits;