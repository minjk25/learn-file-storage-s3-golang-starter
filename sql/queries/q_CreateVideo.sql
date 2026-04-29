-- name: CreateVideo :one
INSERT INTO videos (
    id,
    created_at,
    updated_at,
    title,
    description,
    user_id
)
VALUES (?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, ?, ?, ?)
RETURNING *;