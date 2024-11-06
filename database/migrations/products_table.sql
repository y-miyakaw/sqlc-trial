-- +migrate Up
CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50),
    price INTEGER,
    company_id INTEGER,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- +migrate Down
