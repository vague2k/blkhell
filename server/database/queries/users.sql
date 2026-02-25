-- name: CreateUser :one
INSERT INTO users (id, username, password_hash, role)
VALUES (?, ?, ?, ?)
RETURNING * ;

-- name: GetUserByUsername :one
SELECT * FROM users
WHERE username = ? ;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = ? ;

-- name: UpdateUser :one
UPDATE users
SET
username = ?,
password_hash = ?,
role = ?,
updated_at = CURRENT_TIMESTAMP
WHERE id = ?
RETURNING * ;

-- name: DeleteUserByUsername :exec
DELETE FROM users
WHERE username = ? ;
