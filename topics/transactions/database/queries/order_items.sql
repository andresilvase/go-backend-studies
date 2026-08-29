-- name: CreateOrderItem :exec
INSERT INTO "order_items" (order_id, product_id, quantity) VALUES ($1, $2, $3);