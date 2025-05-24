-- name: GetUserByEmail :one
SELECT id, email, password_hash, role, confirmed, google_id, created_at
FROM users
WHERE email = $1 LIMIT 1;

-- name: CreateUser :exec
INSERT INTO users (id, email, password_hash, role, created_at)
VALUES ($1, $2, $3, $4, $5);

-- name: CheckEmailExists :one
SELECT EXISTS(SELECT 1 FROM users WHERE email = $1);

-- name: UpdateUserConfirmation :exec
UPDATE users SET confirmed = true WHERE id = $1;

-- name: UpdateUserGoogleID :exec
UPDATE users SET google_id = $1 WHERE id = $2; 