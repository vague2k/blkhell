-- name: GetDashboardStats :one
SELECT
    (SELECT COUNT(*) FROM files WHERE owner_type = 'label') AS label_assets,
    (SELECT COUNT(*) FROM bands) AS bands,
    (SELECT COUNT(*) FROM releases) AS releases,
    (SELECT COUNT(*) FROM projects) AS projects;

-- name: GetDashboardBands :many
SELECT
    b.id AS band_id,
    b.name AS band_name,
    b.country,
    b.created_at,
    COUNT(DISTINCT r.id) AS release_count,
    COUNT(DISTINCT p.id) AS project_count,
    COUNT(DISTINCT CASE WHEN p.status = 'done' THEN p.id END) AS projects_done
FROM bands b
LEFT JOIN releases r ON b.id = r.band_id
LEFT JOIN projects p ON b.id = p.band_id
GROUP BY b.id, b.name, b.country, b.created_at
ORDER BY b.created_at DESC;
