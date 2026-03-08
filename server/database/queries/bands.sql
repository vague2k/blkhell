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

-- name: UpdateBand :one
UPDATE bands
SET name = ?,
country = ?,
removed = ?,
updated_at = CURRENT_TIMESTAMP
WHERE id = ?
RETURNING * ;

-- name: GetBandsFromPreviousYear :many
SELECT * FROM bands WHERE created_at >= DATE('now', '-1 year');

-- name: DeleteBand :exec
DELETE FROM bands
WHERE id = ?
