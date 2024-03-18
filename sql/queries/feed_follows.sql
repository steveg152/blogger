-- name: FollowFeed :one

INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;


-- name: DeleteFeedFollow :one
DELETE FROM feed_follows WHERE user_id = $1 AND id = $2
RETURNING *;

-- name: GetFeedFollowsByUserID :many
SELECT * FROM feed_follows WHERE user_id = $1;


