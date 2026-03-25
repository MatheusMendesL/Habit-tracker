-- name: GetUserByID :one
SELECT id, name, email
FROM users
WHERE id = ?;

-- name: SearchUser :many
SELECT id, name, email
FROM users
WHERE
    (name LIKE CONCAT('%', sqlc.arg(name), '%') OR sqlc.arg(name) = '')
  AND
    (email LIKE CONCAT('%', sqlc.arg(email), '%') OR sqlc.arg(email) = '')
    LIMIT 20;

-- name: UpdateUser :exec
UPDATE users
SET name = COALESCE(?, name),
    email = COALESCE(?, email)
WHERE id = ?;

-- name: UpdatePassword :exec
UPDATE users
SET password = ?
WHERE id = ?;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = ?;

-- name: GetUsersByIDs :many
SELECT id, name, email
FROM users
WHERE id IN (sqlc.slice('user_ids'));
