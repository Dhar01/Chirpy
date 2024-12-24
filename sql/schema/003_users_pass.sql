-- +goose Up
ALTER TABLE users
ADD COLUMN hashed_password TEXT NOT NULL DEFAULT 'unset';

-- +goose Down
ALTER TABLE users
DROP TABLE IF EXISTS hashed_password;