-- name: StartFollowing :exec
INSERT INTO follows (follower_id, followee_id)
VALUES ($1, $2);

-- name: Unfollow :exec
DELETE FROM follows
WHERE follower_id = $1
  AND followee_id = $2;

-- name: ListFollowers :many
SELECT follower_id
FROM follows
WHERE followee_id = $1;

-- name: ListFollowing :many
SELECT followee_id
FROM follows
WHERE follower_id = $1;