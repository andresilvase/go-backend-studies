CREATE TABLE orders_items(
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    order_id BIGINT NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    product_id BIGINT NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    quantity INT NOT NULL check (quantity > 0),
    UNIQUE(order_id, product_id)
)

---- create above / drop below ----

DROP TABLE IF EXISTS orders_items;