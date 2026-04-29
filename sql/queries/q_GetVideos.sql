-- name: GetVideos :many
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
WHERE user_id = ?
ORDER BY created_at DESC;