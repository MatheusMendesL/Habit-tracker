-- name: CreateRoutine :one
INSERT INTO routines (user_id, name)
VALUES ($1, $2) RETURNING *;

-- name: GetRoutineByID :one
SELECT id, user_id, name, created_at
FROM routines
WHERE id = $1;

-- name: ListRoutinesByUser :many
SELECT id, user_id, name, created_at
FROM routines
WHERE user_id = $1;

-- name: UpdateRoutine :one
UPDATE routines
SET name = COALESCE(sqlc.narg(name), name)
WHERE id = sqlc.arg(id) RETURNING *;

-- name: DeleteRoutine :exec
DELETE
FROM routines
WHERE id = $1;

-- name: CreateHabit :one
INSERT INTO habits (user_id,
                    name,
                    description,
                    image_url,
                    created_at)
VALUES ($1,
        $2,
        $3,
        $4,
        NOW()) RETURNING *;

-- name: GetHabitByID :one
SELECT id, user_id, name, description, image_url, created_at
FROM habits
WHERE id = $1;

-- name: ListHabitsByUser :many
SELECT id, user_id, name, description, image_url, created_at
FROM habits
WHERE user_id = $1;

-- name: UpdateHabit :one
UPDATE habits
SET name        = COALESCE(sqlc.narg(name), name),
    description = COALESCE(sqlc.narg(description), description),
    image_url   = COALESCE(sqlc.narg(image_url), image_url)
WHERE id = sqlc.arg(id) RETURNING *;

-- name: DeleteHabit :exec
DELETE
FROM habits
WHERE id = $1;

-- name: AddHabitToRoutine :exec
INSERT INTO routine_habits (routine_id, habit_id)
VALUES ($1, $2);

-- name: RemoveHabitFromRoutine :exec
DELETE
FROM routine_habits
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
VALUES ($1, $2) ON CONFLICT (habit_id, completed_at) DO NOTHING;

-- name: UnmarkHabitCompleted :exec
DELETE
FROM habit_logs
WHERE habit_id = $1
  AND completed_at >= DATE (@completed_at:: timestamp)
  AND completed_at
    < DATE (@completed_at:: timestamp) + INTERVAL '1 day';

-- name: GetHabitLogs :many
SELECT habit_id,
       completed_at
FROM habit_logs
WHERE habit_id = $1
  AND completed_at >= @start_date::timestamp
  AND completed_at <= @end_date::timestamp
ORDER BY completed_at DESC;

-- name: MarkRoutineCompleted :exec
INSERT INTO routine_logs (routine_id, completed_at)
VALUES ($1, $2) ON CONFLICT (routine_id, completed_at) DO NOTHING;

-- name: UnmarkRoutineCompleted :exec
DELETE
FROM routine_logs
WHERE routine_id = $1
  AND completed_at >= DATE (@completed_at:: timestamp)
  AND completed_at
    < DATE (@completed_at:: timestamp) + INTERVAL '1 day';

-- name: GetRoutineLogs :many
SELECT routine_id,
       completed_at
FROM routine_logs
WHERE routine_id = $1
  AND completed_at >= @start_date::timestamp
  AND completed_at <= @end_date::timestamp
ORDER BY completed_at DESC;
