-- name: CreateFile :one
INSERT INTO files (
    id,
    user_id,
    path,
    filename,
    ext,
    size,
    mimetype,
    owner_type,
    owner_id,
    audio_duration_seconds
)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
RETURNING * ;

-- name: GetFileByID :one
SELECT * FROM files
WHERE id = ? ;

-- name: GetFiles :many
SELECT * FROM files
ORDER BY filename ;

-- name: GetLabelImageFiles :many
SELECT * FROM files
WHERE owner_type = 'label'
AND mimetype LIKE 'image/%'
ORDER BY filename ;

-- name: GetBandImageFilesByID :many
SELECT * FROM files
WHERE owner_type = 'band'
AND mimetype LIKE 'image/%'
AND owner_id = ?
ORDER BY filename ;

-- name: GetReleaseImageFilesByID :many
SELECT * FROM files
WHERE owner_type = 'release'
AND mimetype LIKE 'image/%'
AND owner_id = ?
ORDER BY filename ;

-- name: GetReleaseAudioFilesByID :many
SELECT * FROM files
WHERE owner_type = 'release'
AND mimetype LIKE 'audio/%'
AND owner_id = ?
ORDER BY filename ;

-- name: GetProjectImageFilesByID :many
SELECT * FROM files
WHERE owner_type = 'project'
AND mimetype LIKE 'image/%'
AND owner_id = ?
ORDER BY filename ;

-- name: GetFileByPartialName :many
SELECT *
FROM files
WHERE filename LIKE ?
OR ext like ? ;

-- name: DeleteFile :one
DELETE FROM files
WHERE id = ?
RETURNING * ;
