-- name: UpdateVideo :exec
UPDATE videos
SET
    updated_at = CURRENT_TIMESTAMP,
    title = ?,
    description = ?,
    thumbnail_url = ?,
    video_url = ?,
    user_id = ?
WHERE id = ?;