-- name: GetDashboardStats :one
SELECT
    (SELECT COUNT(*) FROM files WHERE owner_type = 'label') AS label_assets,
    (SELECT COUNT(*) FROM bands) AS bands,
    (SELECT COUNT(*) FROM releases) AS releases,
    (SELECT COUNT(*) FROM projects) AS projects;
