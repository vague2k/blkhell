-- name: CreateBand :one
INSERT INTO bands (
    id,
    name,
    country,
    removed
) VALUES (?, ?, ?, ?)
RETURNING * ;

-- name: GetBandByID :one
SELECT * FROM bands
WHERE id = ? ;

-- name: GetBands :many
SELECT * FROM bands
ORDER BY name ;

-- name: GetActiveBands :many
SELECT * FROM bands
WHERE removed = 0
ORDER BY name ;

-- name: UpdateBand :one
UPDATE bands
SET name = ?,
country = ?,
updated_at = CURRENT_TIMESTAMP
WHERE id = ?
RETURNING * ;

-- name: GetBandsFromPreviousYear :many
SELECT * FROM bands WHERE created_at >= DATE('now', '-1 year');

-- name: RemoveBand :one
UPDATE bands
SET removed = 1,
updated_at = CURRENT_TIMESTAMP
WHERE id = ?
RETURNING * ;

-- name: DeleteBand :exec
DELETE FROM bands WHERE id = ?;
