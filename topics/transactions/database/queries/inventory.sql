-- name: SetProductInventory :exec
INSERT INTO inventory(product_id, stock) VALUES ($1, $2);

-- name: GetStockForProductID :one
SELECT stock FROM inventory WHERE product_id = $1;

-- name: UpdateProductInventory :exec
UPDATE "inventory" SET stock = $1 WHERE product_id = $2;