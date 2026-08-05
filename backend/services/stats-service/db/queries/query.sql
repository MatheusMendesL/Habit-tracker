-- name: CreateUserStats :one
INSERT INTO user_stats(user_id)
VALUES ($1) RETURNING *;

-- name: GetUserStats :one
SELECT * FROM user_stats 
WHERE user_id = $1;

-- name: DeleteUserStats :exec
DELETE 
FROM user_stats 
WHERE user_id = $1;

-- name: IncrementCompletedHabits :exec
UPDATE user_stats
SET completed_habits = completed_habits + 1,
    updated_at = NOW()
WHERE user_id = sqlc.arg(id);

-- name: DecrementCompletedHabits :exec
UPDATE user_stats
SET completed_habits = GREATEST(completed_habits - 1, 0),
    updated_at = NOW()
WHERE user_id = sqlc.arg(id);

-- name: IncrementCompletedRoutines :exec
UPDATE user_stats
SET completed_routines = completed_routines + 1,
    updated_at = NOW()
WHERE user_id = sqlc.arg(id);

-- name: DecrementCompletedRoutines :exec
UPDATE user_stats
SET completed_routines = GREATEST(completed_routines - 1, 0),
    updated_at = NOW()
WHERE user_id = sqlc.arg(id);

-- name: UpdateHabitStreaks :exec
UPDATE user_stats
SET current_habit_streak = sqlc.arg(current_habit_streak),
    longest_habit_streak = sqlc.arg(longest_habit_streak),
    updated_at = NOW()
WHERE user_id = sqlc.arg(id);

-- name: UpdateRoutineStreaks :exec
UPDATE user_stats
SET longest_routine_streak = sqlc.arg(longest_routine_streak),
    current_routine_streak = sqlc.arg(current_routine_streak),
    updated_at = NOW()
WHERE user_id = sqlc.arg(id);