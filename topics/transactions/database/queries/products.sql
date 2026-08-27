-- name: CreateProduct :exec
INSERT INTO products(name, price) VALUES ($1, $2);