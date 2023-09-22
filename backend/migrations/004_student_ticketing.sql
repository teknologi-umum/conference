-- +goose Up
ALTER TABLE ticketing
    ADD COLUMN IF NOT EXISTS student BOOLEAN DEFAULT FALSE;

-- +goose Down
ALTER TABLE ticketing
    DROP COLUMN student;