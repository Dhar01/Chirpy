// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: chirps.sql

package database

import (
	"context"

	"github.com/google/uuid"
)

const createChirp = `-- name: CreateChirp :one
INSERT INTO chirps (id, created_at, updated_at, body, user_id)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING id, created_at, updated_at, body, user_id
`

type CreateChirpParams struct {
	Body   string
	UserID uuid.UUID
}

func (q *Queries) CreateChirp(ctx context.Context, arg CreateChirpParams) (Chirp, error) {
	row := q.db.QueryRowContext(ctx, createChirp, arg.Body, arg.UserID)
	var i Chirp
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Body,
		&i.UserID,
	)
	return i, err
}

const deleteChirp = `-- name: DeleteChirp :exec
DELETE FROM chirps
WHERE id = $1 AND user_id = $2
`

type DeleteChirpParams struct {
	ID     uuid.UUID
	UserID uuid.UUID
}

func (q *Queries) DeleteChirp(ctx context.Context, arg DeleteChirpParams) error {
	_, err := q.db.ExecContext(ctx, deleteChirp, arg.ID, arg.UserID)
	return err
}

const getAllChirps = `-- name: GetAllChirps :many
SELECT id, created_at, updated_at, body, user_id FROM chirps ORDER BY created_at ASC
`

func (q *Queries) GetAllChirps(ctx context.Context) ([]Chirp, error) {
	rows, err := q.db.QueryContext(ctx, getAllChirps)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Chirp
	for rows.Next() {
		var i Chirp
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Body,
			&i.UserID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getIDofChirp = `-- name: GetIDofChirp :one
SELECT user_id FROM chirps
WHERE id = $1
`

func (q *Queries) GetIDofChirp(ctx context.Context, id uuid.UUID) (uuid.UUID, error) {
	row := q.db.QueryRowContext(ctx, getIDofChirp, id)
	var user_id uuid.UUID
	err := row.Scan(&user_id)
	return user_id, err
}

const getSingleChirp = `-- name: GetSingleChirp :one
SELECT id, created_at, updated_at, body, user_id FROM chirps
WHERE id = $1
`

func (q *Queries) GetSingleChirp(ctx context.Context, id uuid.UUID) (Chirp, error) {
	row := q.db.QueryRowContext(ctx, getSingleChirp, id)
	var i Chirp
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Body,
		&i.UserID,
	)
	return i, err
}
