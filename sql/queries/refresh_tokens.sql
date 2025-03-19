-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens(token, created_at, updated_at, user_id, expires_at, revoked_at)
VALUES (
    $1,
    NOW(),
    NOW(),
    $2,
    CURRENT_TIMESTAMP + INTERVAL '60 days',
    NULL
)
RETURNING *;

-- name: GetRefreshToken :one
SELECT token, user_id, expires_at, revoked_at FROM refresh_tokens
WHERE token = $1;

-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens
SET revoked_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
WHERE token = $1 AND revoked_at IS NULL;