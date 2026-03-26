-- name: StartFollowing :exec
INSERT INTO follows (follower_id, followee_id)
VALUES (?, ?);

-- name: Unfollow :exec
DELETE FROM follows
WHERE follower_id = ? AND followee_id = ?;

-- name: ListFollowers :many
SELECT follower_id
FROM follows
WHERE followee_id = ?;

-- name: ListFollowing :many
SELECT followee_id
FROM follows
WHERE follower_id = ?;