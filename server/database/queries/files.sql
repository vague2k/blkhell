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
    owner_id
)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
RETURNING * ;

-- name: GetFileByID :one
SELECT * FROM files
WHERE id = ? ;

-- name: GetFiles :many
SELECT * FROM files
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
