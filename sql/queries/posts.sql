-- name: CreatePost :one

INSERT INTO posts (id, created_at, updated_at, feed_id, title, url, published_at, description)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;


-- name: GetPostByURL :one
SELECT * FROM posts WHERE url = $1;

-- name: GetPostsByFeedID :many
SELECT * FROM posts WHERE feed_id = $1;

-- name: GetPostByID :one
SELECT * FROM posts WHERE id = $1;


-- name: GetPostsByUserId :many
SELECT p.* FROM posts p
JOIN feeds f ON p.feed_id = f.id
WHERE f.user_id = $1 LIMIT $2 OFFSET $3;

