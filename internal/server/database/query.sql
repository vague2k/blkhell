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

-- name: CreateSession :one
INSERT INTO sessions (id, token, user_id, expires_at)
VALUES (?, ?, ?, ?)
RETURNING * ;

-- name: GetSessionByToken :one
SELECT * FROM sessions
WHERE token = ?
AND expires_at > CURRENT_TIMESTAMP ;

-- name: UpdateUser :one
UPDATE users
SET username = ?, password_hash = ?, role = ?
WHERE id = ?
RETURNING * ;

-- name: DeleteSession :exec
DELETE FROM sessions
WHERE id = ? ;

-- name: DeleteExpiredSessions :exec
DELETE FROM sessions
WHERE expires_at < CURRENT_TIMESTAMP ;

-- name: DeleteUserByUsername :exec
DELETE FROM users WHERE username = ? ;

-- name: CreateImage :one
INSERT INTO images (id, user_id, path, filename, ext, size)
VALUES (?, ?, ?, ?, ?, ?)
RETURNING * ;

-- name: GetImages :many
SELECT * FROM images
ORDER BY filename ;
