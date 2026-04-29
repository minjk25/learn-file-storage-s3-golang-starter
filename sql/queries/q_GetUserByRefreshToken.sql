-- name: GetUserByRefreshToken :one
SELECT u.id, u.email, u.created_at, u.updated_at, u.password
FROM users u
JOIN refresh_tokens rt ON u.id = rt.user_id
WHERE rt.token = ?;