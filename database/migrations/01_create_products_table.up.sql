-- +migrate Up
CREATE TABLE IF NOT EXISTS products (
  id VARCHAR(10) PRIMARY KEY,
  user_id VARCHAR(10),
  name VARCHAR(50) NOT NULL,
  price DECIMAL(10, 2) NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- +migrate Down