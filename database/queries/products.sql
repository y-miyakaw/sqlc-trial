-- name: GetProduct :one
SELECT * FROM products WHERE id = $1;

-- name: GetAllProducts :many
SELECT * FROM products;

-- name: GetProductsByUserIDAndColor :many
SELECT * FROM products
WHERE (user_id = COALESCE($1, user_id))
AND (color = COALESCE($2, color))
AND (name =  COALESCE(NULLIF($3, ''), name));

-- name: GetProductsByIDsAndColor :many
SELECT * FROM products
WHERE (COALESCE(NULLIF($1, '')::text[], ARRAY[]::text[]) = ARRAY[]::text[] OR id = ANY($1::text[]))
AND (color = COALESCE($2, color));

-- name: CreateProduct :one
INSERT INTO products (id, user_id, name, price, identifier, color) VALUES ($1, $2, $3, $4, $5, $6) RETURNING *;

-- name: UpdateProduct :one
UPDATE products SET user_id = $2, name = $3, price = $4 WHERE id = $1 RETURNING *;

-- name: DeleteProduct :one
DELETE FROM products WHERE id = $1 RETURNING *;