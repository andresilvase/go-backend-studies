.PHONY: migrations rollback-all sql db_fed reset status

-include topics/transactions/database/migrations/.env
export

migrations:
	go run ./cmd/terndotenv

rollback-all:
	tern migrate -m topics/transactions/database/migrations \
		-c topics/transactions/database/migrations/tern.conf -d 0

sql:
	sqlc generate -f ./topics/transactions/database/sqlc.yml
	
db_fed:
	go run ./topics/transactions/database/db_fed

reset:
	$(MAKE) rollback-all && $(MAKE) migrations && $(MAKE) db_fed

status:
	tern status -m topics/transactions/database/migrations -c topics/transactions/database/migrations/tern.conf