-- name: CreateWallet :one
INSERT INTO wallets(user_id, balance) VALUES ($1, $2) RETURNING id;