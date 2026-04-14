-- name: CreateRoutine :execresult
INSERT INTO routines (user_id, name)
VALUES (?, ?);

-- name: GetRoutineByID :one
SELECT id, user_id, name, created_at
FROM routines
WHERE id = ?;

-- name: ListRoutinesByUser :many
SELECT id, user_id, name, created_at
FROM routines
WHERE user_id = ?;

-- name: UpdateRoutine :exec
UPDATE routines
SET name = COALESCE(?, name)
WHERE id = ?;

-- name: DeleteRoutine :exec
DELETE FROM routines
WHERE id = ?;

-- name: CreateHabit :execresult
INSERT INTO habits (user_id, name, description, image_url)
VALUES (?, ?, ?, ?);

-- name: GetHabitByID :one
SELECT id, user_id, name, description, image_url, created_at
FROM habits
WHERE id = ?;

-- name: ListHabitsByUser :many
SELECT id, user_id, name, description, image_url, created_at
FROM habits
WHERE user_id = ?;

-- name: UpdateHabit :exec
UPDATE habits
SET name        = COALESCE(?, name),
    description = COALESCE(?, description),
    image_url   = COALESCE(?, image_url)
WHERE id = ?;

-- name: DeleteHabit :exec
DELETE FROM habits
WHERE id = ?;

-- name: AddHabitToRoutine :exec
INSERT INTO routine_habits (routine_id, habit_id)
VALUES (?, ?);

-- name: RemoveHabitFromRoutine :exec
DELETE FROM routine_habits
WHERE routine_id = ? AND habit_id = ?;

-- name: ListHabitsByRoutine :many
SELECT h.id, h.user_id, h.name, h.description, h.image_url, h.created_at
FROM habits h
         JOIN routine_habits rh ON rh.habit_id = h.id
WHERE rh.routine_id = ?;

-- name: ListRoutinesByHabit :many
SELECT r.id, r.user_id, r.name, r.created_at
FROM routines r
         JOIN routine_habits rh ON rh.routine_id = r.id
WHERE rh.habit_id = ?;

-- name: MarkHabitCompleted :exec
INSERT INTO habit_logs (habit_id, completed_at)
VALUES (?, ?)
    ON DUPLICATE KEY UPDATE completed_at = completed_at;

-- name: UnmarkHabitCompleted :exec
DELETE FROM habit_logs
WHERE habit_id = ?
  AND DATE(completed_at) = DATE(?);

-- name: GetHabitLogs :many
SELECT habit_id, completed_at
FROM habit_logs
WHERE habit_id = ?
  AND completed_at BETWEEN ? AND ?
ORDER BY completed_at DESC;