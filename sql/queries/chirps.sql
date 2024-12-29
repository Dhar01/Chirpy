-- name: CreateChirp :one
INSERT INTO chirps (id, created_at, updated_at, body, user_id)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING *;

-- name: GetAllChirpsASC :many
SELECT * FROM chirps
WHERE (user_id = $1 OR $1 = '00000000-0000-0000-0000-000000000000')
ORDER BY created_at ASC;

-- name: GetAllChirpsDESC :many
SELECT * FROM chirps
WHERE (user_id = $1 OR $1 = '00000000-0000-0000-0000-000000000000')
ORDER BY created_at DESC;

-- name: GetSingleChirp :one
SELECT * FROM chirps
WHERE id = $1;

-- name: DeleteChirp :exec
DELETE FROM chirps
WHERE id = $1 AND user_id = $2;

-- name: GetIDofChirp :one
SELECT user_id FROM chirps
WHERE id = $1;