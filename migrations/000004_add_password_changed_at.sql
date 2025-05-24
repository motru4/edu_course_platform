-- +goose Up
ALTER TABLE users ADD COLUMN IF NOT EXISTS password_changed_at TIMESTAMPTZ DEFAULT NOW();

-- +goose Down
ALTER TABLE users DROP COLUMN IF EXISTS password_changed_at; 