-- name: GetProduct :one
SELECT *
FROM products
WHERE id = sqlc.arg(id);

-- name: GetAllCompanies :many
SELECT *
FROM company;

-- name: GetAllProducts :many
SELECT *
FROM products;

-- name: GetProductsByNameOrPriceOrCompanyID :many
SELECT *
FROM products
WHERE 1=1
AND (
    CASE
        WHEN sqlc.narg(name)::text IS NULL THEN TRUE
        ELSE name = sqlc.narg(name)
    END
)
AND (
    CASE
        WHEN sqlc.narg(price)::integer IS NULL THEN TRUE
        ELSE price = sqlc.narg(price)
    END
)
AND (
    CASE
        WHEN sqlc.narg(company_id)::integer IS NULL THEN TRUE
        ELSE name = sqlc.narg(company_id)
    END
);


-- name: GetProductsAndCompanyByCompanyID :many
SELECT *
FROM products
JOIN company ON company.id = products.company_id
WHERE company.id = sqlc.arg(id);

-- name: CreateProductWithReturning :one
INSERT INTO products (name, price, company_id)
VALUES (sqlc.narg(name), sqlc.narg(price), sqlc.narg(company_id))
RETURNING *;

-- name: CreateProductWithoutReturning :exec
INSERT INTO products (name, price, company_id)
VALUES (sqlc.narg(name), sqlc.narg(price), sqlc.narg(company_id));

-- name: CreateCompanyWithoutReturning :exec
INSERT INTO company (name, address, person)
VALUES (sqlc.arg(name), sqlc.narg(address), sqlc.narg(person));

-- name: UpdateProduct :exec
UPDATE products
SET
    name = sqlc.narg(name),
    price = sqlc.narg(price),
    company_id = sqlc.narg(company_id)
WHERE id = sqlc.arg(id);

-- name: DeleteProduct :exec
DELETE
FROM products
WHERE id = sqlc.arg(id);