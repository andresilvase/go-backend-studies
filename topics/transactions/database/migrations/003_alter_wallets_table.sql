ALTER TABLE wallets
ADD CONSTRAINT check_positive_balance 
CHECK (balance >= 0);

---- create above / drop below ----

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
