-- name: CreateRoutine :execresult
INSERT INTO routines (user_id, name)
VALUES ($1, $2);

-- name: GetRoutineByID :one
SELECT id, user_id, name, created_at
FROM routines
WHERE id = $1;

-- name: ListRoutinesByUser :many
SELECT id, user_id, name, created_at
FROM routines
WHERE user_id = $1;

-- name: UpdateRoutine :exec
UPDATE routines
SET name = COALESCE($1, name)
WHERE id = $2;

-- name: DeleteRoutine :exec
DELETE FROM routines
WHERE id = $1;

-- name: CreateHabit :execresult
INSERT INTO habits (user_id, name, description, image_url)
VALUES ($1, $2, $3, $4);

-- name: GetHabitByID :one
SELECT id, user_id, name, description, image_url, created_at
FROM habits
WHERE id = $1;

-- name: ListHabitsByUser :many
SELECT id, user_id, name, description, image_url, created_at
FROM habits
WHERE user_id = $1;

-- name: UpdateHabit :exec
UPDATE habits
SET name        = COALESCE($1, name),
    description = COALESCE($2, description),
    image_url   = COALESCE($3, image_url)
WHERE id = $4;

-- name: DeleteHabit :exec
DELETE FROM habits
WHERE id = $1;

-- name: AddHabitToRoutine :exec
INSERT INTO routine_habits (routine_id, habit_id)
VALUES ($1, $2);

-- name: RemoveHabitFromRoutine :exec
DELETE FROM routine_habits
WHERE routine_id = $1
  AND habit_id = $2;

-- name: ListHabitsByRoutine :many
SELECT h.id,
       h.user_id,
       h.name,
       h.description,
       h.image_url,
       h.created_at
FROM habits h
         JOIN routine_habits rh ON rh.habit_id = h.id
WHERE rh.routine_id = $1;

-- name: ListRoutinesByHabit :many
SELECT r.id,
       r.user_id,
       r.name,
       r.created_at
FROM routines r
         JOIN routine_habits rh ON rh.routine_id = r.id
WHERE rh.habit_id = $1;

-- name: MarkHabitCompleted :exec
INSERT INTO habit_logs (habit_id, completed_at)
VALUES ($1, $2)
    ON CONFLICT (habit_id, completed_at) DO NOTHING;

-- name: UnmarkHabitCompleted :exec
DELETE FROM habit_logs
WHERE habit_id = $1
  AND DATE(completed_at) = DATE($2);

-- name: GetHabitLogs :many
SELECT habit_id, completed_at
FROM habit_logs
WHERE habit_id = $1
  AND completed_at BETWEEN $2 AND $3
ORDER BY completed_at DESC;