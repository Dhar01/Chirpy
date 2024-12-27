-- name: CreateRefreshToken :exec
INSERT INTO refresh_tokens (
    refreshToken,
    user_id,
    created_at,
    updated_at,
    expires_at,
    revoked_at
) VALUES (
    $1,
    $2,
    NOW(),
    NOW(),
    NOW() + INTERVAL '60 days',
    NULL
) RETURNING *;

-- name: GetUserFromRefreshToken :one
SELECT user_id FROM refresh_tokens
WHERE refreshToken = $1
    AND expires_at > NOW()
    AND revoked_at IS NULL;
