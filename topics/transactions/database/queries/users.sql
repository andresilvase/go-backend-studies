-- name: CreateUser :exec
INSERT INTO users(name) VALUES ($1);