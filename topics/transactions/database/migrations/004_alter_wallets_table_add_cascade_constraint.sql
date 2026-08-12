ALTER TABLE wallets
    DROP CONSTRAINT wallets_user_id_fkey,
    ADD CONSTRAINT wallets_user_id_fkey
    FOREIGN KEY (user_id)
    REFERENCES users(id)
    ON DELETE CASCADE;




---- create above / drop below ----

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
