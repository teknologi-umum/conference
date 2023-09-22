-- +goose Up
CREATE TABLE ticketing
(
    id                 UUID PRIMARY KEY,
    email              VARCHAR(255) UNIQUE NOT NULL,
    receipt_photo_path TEXT                NOT NULL,
    paid               BOOLEAN   DEFAULT FALSE,
    sha256sum          BYTEA               NULL,
    used               BOOLEAN   DEFAULT FALSE,
    created_at         TIMESTAMP DEFAULT NOW(),
    updated_at         TIMESTAMP DEFAULT NOW()
);

-- +goose Down
DROP TABLE ticketing;