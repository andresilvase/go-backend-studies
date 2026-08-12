CREATE TABLE wallets (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL UNIQUE REFERENCES users(id),
    balance BIGINT NOT NULL DEFAULT 0
);

---- create above / drop below ----

DROP TABLE IF EXISTS wallets;

