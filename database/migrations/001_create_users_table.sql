CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL
);

---- create above / drop below ----

DROP TABLE IF EXISTS users CASCADE;