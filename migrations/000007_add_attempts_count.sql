-- +goose Up
ALTER TABLE lesson_progress
ADD COLUMN attempts_count INTEGER NOT NULL DEFAULT 0;

-- +goose Down
ALTER TABLE lesson_progress
DROP COLUMN attempts_count; 