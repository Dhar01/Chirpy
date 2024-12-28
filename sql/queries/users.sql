-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE email = $1;

-- name: UpdateEmail :exec
UPDATE users
SET
    email = $1,
    updated_at = NOW()
WHERE id = $2;

-- name: UpdatePassword :exec
UPDATE users
SET
    hashed_password = $1,
    updated_at = NOW()
WHERE id = $2;

-- name: Reset :exec
DELETE FROM users;