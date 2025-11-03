-- +goose Up
ALTER TABLE users ADD COLUMN IF NOT EXISTS phone VARCHAR(20);

-- +goose Down
ALTER TABLE users DROP COLUMN IF EXISTS phone;
