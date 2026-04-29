-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (token, created_at, updated_at, user_id, expires_at) 
VALUES (?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, ?, ?)
RETURNING *;