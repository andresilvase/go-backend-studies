-- name: CreateOrder :one
INSERT INTO "orders"(user_id) VALUES ($1) RETURNING id;