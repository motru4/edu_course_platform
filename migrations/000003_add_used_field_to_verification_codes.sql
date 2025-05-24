-- +goose Up
ALTER TABLE verification_codes
ADD COLUMN used BOOLEAN NOT NULL DEFAULT FALSE;

-- +goose Down
ALTER TABLE verification_codes
DROP COLUMN used;
