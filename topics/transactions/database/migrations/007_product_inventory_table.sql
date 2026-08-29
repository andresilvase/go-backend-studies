CREATE TABLE inventory(
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    product_id BIGINT NOT NULL UNIQUE REFERENCES products(id) ON DELETE CASCADE,
    stock BIGINT NOT NULL check (stock>=0)
)

---- create above / drop below ----

DROP TABLE IF EXISTS inventory;