-- name: GetUser :one
SELECT id, created_at, updated_at, email, password
FROM users
WHERE id = ?;