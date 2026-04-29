-- name: DeleteVideo :exec
DELETE FROM videos
WHERE id = ?;