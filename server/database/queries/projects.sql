-- name: CreateProject :one
INSERT INTO projects (
    id,
    band_id,
    release_id,
    name,
    type
) VALUES (?, ?, ?, ?, ?)
RETURNING * ;

-- name: GetProjects :many
SELECT * FROM projects
ORDER BY name ;

-- name: GetProjectByID :one
SELECT *
FROM projects
WHERE id = ? ;

-- name: GetProjectsByRelease :many
SELECT *
FROM projects
WHERE release_id = ?
ORDER BY created_at DESC ;

-- name: UpdateProject :one
UPDATE projects
SET name = ?,
type = ?,
updated_at = CURRENT_TIMESTAMP
WHERE id = ?
RETURNING * ;

-- name: DeleteProject :exec
DELETE FROM projects
WHERE id = ? ;
