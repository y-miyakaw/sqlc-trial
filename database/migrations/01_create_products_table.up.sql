-- +migrate Up
CREATE TABLE products (
  id VARCHAR(10) PRIMARY KEY,
  user_id VARCHAR(10),
  name VARCHAR(50) NOT NULL,
  price DECIMAL(10, 2) NOT NULL,
  identifier VARCHAR(10),
  color VARCHAR(10),
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- +migrate Down