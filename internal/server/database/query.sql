-- name: CreateUser :one
INSERT INTO users (username, password_hash, role)
VALUES (?, ?, ?)
RETURNING * ;

-- name: GetUserByUsername :one
SELECT * FROM users
WHERE username = ? ;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = ? ;

-- name: CreateSession :one
INSERT INTO sessions (id, user_id, expires_at)
VALUES (?, ?, ?)
RETURNING * ;

-- name: GetSession :one
SELECT * FROM sessions
WHERE id = ? ;

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
