-- name: GetProduct :one
SELECT * FROM products WHERE id = $1;

-- name: GetAllProducts :many
SELECT * FROM products;

-- name: CreateProduct :one
INSERT INTO products (user_id, name, price) VALUES ($1, $2, $3) RETURNING *;

-- name: UpdateProduct :one
UPDATE products SET user_id = $2, name = $3, price = $4 WHERE id = $1 RETURNING *;

-- name: DeleteProduct :one
DELETE FROM products WHERE id = $1 RETURNING *;