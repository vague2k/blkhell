-- name: CreateRelease :one
INSERT INTO releases (
    id,
    band_id,
    name,
    type,
    number
) VALUES (?, ?, ?, ?, ?)
RETURNING * ;

-- name: GetReleases :many
SELECT * FROM releases
ORDER BY name ;

-- name: GetReleaseByID :one
SELECT * FROM releases
WHERE id = ? ;

-- name: GetReleasesByBand :many
SELECT * FROM releases
WHERE band_id = ?
ORDER BY name DESC ;

-- name: GetReleasesFromPreviousYear :many
SELECT * FROM releases WHERE created_at >= DATE('now', '-1 year');

-- name: UpdateRelease :one
UPDATE releases
SET name = ?,
type = ?,
number = ?,
updated_at = CURRENT_TIMESTAMP
WHERE id = ?
RETURNING * ;

-- name: DeleteRelease :exec
DELETE FROM releases
WHERE id = ? ;
