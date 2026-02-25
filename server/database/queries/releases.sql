-- name: CreateRelease :one
INSERT INTO releases (
    id,
    band_id,
    name,
    type,
    number
) VALUES (?, ?, ?, ?, ?)
RETURNING * ;

-- name: GetReleaseByID :one
SELECT * FROM releases
WHERE id = ? ;

-- name: GetReleasesByBand :many
SELECT * FROM releases
WHERE band_id = ?
ORDER BY name DESC ;

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
