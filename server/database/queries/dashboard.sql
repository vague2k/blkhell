-- name: GetDashboardStats :one
SELECT
    (SELECT COUNT(*) FROM files WHERE owner_type = 'label')
        AS label_assets_count,
    (SELECT COUNT(*) FROM bands) AS bands_count,
    (SELECT COUNT(*) FROM releases) AS releases_count,
    (SELECT COUNT(*) FROM projects) AS projects_count,
    (SELECT id FROM bands ORDER BY created_at DESC LIMIT 1) AS latest_band_id,
    (SELECT name FROM bands ORDER BY created_at DESC LIMIT 1
    ) AS latest_band_name,
    (SELECT id FROM releases ORDER BY created_at DESC LIMIT 1
    ) AS latest_release_id,
    (SELECT name FROM releases ORDER BY created_at DESC LIMIT 1
    ) AS latest_release_title,
    (SELECT id FROM projects ORDER BY created_at DESC LIMIT 1
    ) AS latest_project_id,
    (SELECT name FROM projects ORDER BY created_at DESC LIMIT 1) AS latest_project_name;

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
