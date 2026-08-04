CREATE TABLE user_stats {
    user_id UUID PRIMARY KEY,
    completed_habits INTEGER NOT NULL DEFAULT 0,
    completed_routines INTEGER NOT NULL DEFAULT 0,
    current_habit_streak INTEGER NOT NULL DEFAULT 0,
    longest_habit_streak INTEGER NOT NULL DEFAULT 0,
    current_routine_streak INTEGER NOT NULL DEFAULT 0,
    longest_routine_streak INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
}