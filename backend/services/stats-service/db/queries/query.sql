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
SET completed_habits = completed_habits + 1
WHERE user_id = sqlc.arg(id);

-- name: DecrementCompletedHabits :exec
UPDATE user_stats
SET completed_habits = completed_habits - 1
WHERE user_id = sqlc.arg(id);

-- name: IncrementCompletedRoutines:exec
UPDATE user_stats
SET completed_routines = completed_routines + 1
WHERE user_id = sqlc.arg(id);

-- name: DecrementCompletedRoutines :exec
UPDATE user_stats
SET completed_routines = completed_routines - 1
WHERE user_id = sqlc.arg(id);

-- name: UpdateHabitStreaks :exec
UPDATE user_stats
SET longest_habit_streak = longest_habit_streak + 1,
    current_habit_streak = current_habit_streak + 1,
WHERE user_id = sqlc.arg(id);

-- name: UpdateRoutinesStreaks :exec
UPDATE user_stats
SET longest_routine_streak = longest_routine_streak + 1,
    current_routine_streak = current_routine_streak + 1,
WHERE user_id = sqlc.arg(id);