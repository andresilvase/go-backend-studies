-- name: CreateProduct :one
INSERT INTO products(name, price) VALUES ($1, $2) RETURNING id;