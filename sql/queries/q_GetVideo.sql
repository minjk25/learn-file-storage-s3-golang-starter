-- name: GetVideo :one
SELECT
    id,
    created_at,
    updated_at,
    title,
    description,
    thumbnail_url,
    video_url,
    user_id
FROM videos
WHERE id = ?;