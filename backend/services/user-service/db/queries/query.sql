-- name: GetUserByID :one
SELECT id, name, email
FROM users
WHERE id = ?;

-- name: SearchUser :many
SELECT id, name, email
FROM users
WHERE ( ? IS NULL OR name LIKE ? )
  AND ( ? IS NULL OR email LIKE ? );

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

-- name: StartFollowing :exec
INSERT INTO followers (follower_id, followee_id)
VALUES (?, ?);

-- name: Unfollow :exec
DELETE FROM followers
WHERE follower_id = ? AND followee_id = ?;

-- name: ListFollowers :many
SELECT u.id, u.name, u.email
FROM followers f
         JOIN users u ON u.id = f.follower_id
WHERE f.followee_id = ?;

-- name: ListFollowing :many
SELECT u.id, u.name, u.email
FROM followers f
         JOIN users u ON u.id = f.followee_id
WHERE f.follower_id = ?;