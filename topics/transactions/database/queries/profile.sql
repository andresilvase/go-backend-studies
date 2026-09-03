-- name: CreateProfile :exec
INSERT INTO "user_profile"(user_id, wallet_id, category) VALUES ($1, $2, $3);