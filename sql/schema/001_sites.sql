-- +goose Up
CREATE TABLE sites(
    id SERIAL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL,
    url TEXT UNIQUE NOT NULL,
    description TEXT,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    status TEXT DEFAULT 'unknown'
);

-- +goose Down
DROP TABLE IF EXISTS sites;