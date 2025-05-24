-- name: CreateRefreshSession :exec
INSERT INTO refresh_sessions (id, user_id, refresh_token, expires_at, created_at)
VALUES ($1, $2, $3, $4, $5);

-- name: GetRefreshSession :one
SELECT id, user_id, refresh_token, expires_at, created_at
FROM refresh_sessions
WHERE refresh_token = $1 AND expires_at > $2;

-- name: DeleteRefreshSession :exec
DELETE FROM refresh_sessions WHERE refresh_token = $1;

-- name: DeleteExpiredSessions :exec
DELETE FROM refresh_sessions WHERE expires_at < $1;

-- name: DeleteUserSessions :exec
DELETE FROM refresh_sessions WHERE user_id = $1; 