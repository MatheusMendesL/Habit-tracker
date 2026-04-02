-- name: CreateHabit :execresult
INSERT INTO habits (user_id, name, description, image_url)
VALUES (?, ?, ?, ?);

-- name: GetHabitByID :one
SELECT id, user_id, name, description, image_url
FROM habits
WHERE id = ?;

-- name: ListHabitsByUser :many
SELECT id, user_id, name, description, image_url
FROM habits
WHERE user_id = ?;

-- name: UpdateHabit :exec
UPDATE habits
SET name        = COALESCE(?, name),
    description = COALESCE(?, description),
    image_url   = COALESCE(?, image_url)
WHERE id = ?;

-- name: DeleteHabit :exec
DELETE
FROM habits
WHERE id = ?;

-- name: CreateRoutine :execresult
INSERT INTO routines (user_id, name)
VALUES (?, ?);

-- name: ListRoutinesByUser :many
SELECT id, user_id, name
FROM routines
WHERE user_id = ?;

-- name: MarkHabitCompleted :exec
INSERT INTO habit_logs (habit_id, completed_at)
VALUES (?, ?) ON DUPLICATE KEY
UPDATE completed_at = completed_at;

-- name: UnmarkHabitCompleted :exec
DELETE
FROM habit_logs
WHERE habit_id = ? AND DATE (completed_at) = DATE (?);

-- name: GetHabitLogs :many
SELECT habit_id, completed_at
FROM habit_logs
WHERE habit_id = ?
  AND completed_at BETWEEN ? AND ?
ORDER BY completed_at DESC;