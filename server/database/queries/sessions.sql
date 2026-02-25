-- name: CreateSession :one
INSERT INTO sessions (id, token, user_id, expires_at)
VALUES (?, ?, ?, ?)
RETURNING * ;

-- name: GetSessionByToken :one
SELECT * FROM sessions
WHERE token = ?
AND expires_at > CURRENT_TIMESTAMP ;

-- name: DeleteSession :exec
DELETE FROM sessions
WHERE token = ? ;

-- name: DeleteExpiredSessions :exec
DELETE FROM sessions
WHERE expires_at < CURRENT_TIMESTAMP ;
