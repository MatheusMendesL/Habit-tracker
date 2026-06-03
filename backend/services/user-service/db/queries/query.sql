-- name: GetUserByID :one
SELECT id, name, email
FROM users
WHERE id = $1;

-- name: SearchUser :many
SELECT id, name, email
FROM users
WHERE
    (name ILIKE '%' || sqlc.arg(name) || '%' OR sqlc.arg(name) = '')
  AND
    (email ILIKE '%' || sqlc.arg(email) || '%' OR sqlc.arg(email) = '')
    LIMIT 20;

-- name: UpdateUser :exec
UPDATE users
SET
    name = COALESCE($1, name),
    email = COALESCE($2, email)
WHERE id = $3;

-- name: UpdatePassword :exec
UPDATE users
SET password = $1
WHERE id = $2;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;

-- name: GetUsersByIDs :many
SELECT id, name, email
FROM users
WHERE id = ANY($1::uuid[]);