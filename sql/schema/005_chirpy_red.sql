-- +goose Up
ALTER TABLE users
ADD COLUMN is_chirpy_red BOOL DEFAULT false;

-- +goose Down
ALTER TABLE users
DROP COLUMN IF EXISTS is_chirpy_red;