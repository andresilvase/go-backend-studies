-include topics/transactions/database/migrations/.env
export

migrations:
	go run ./cmd/terndotenv

rollback-all:
	tern migrate -m topics/transactions/database/migrations \
		-c topics/transactions/database/migrations/tern.conf -d 0

sql:
	sqlc generate -f ./topics/transactions/database/sqlc.yml
	
backfill:
	go run ./topics/transactions/database/migrations

status:
	tern status -m topics/transactions/database/migrations -c topics/transactions/database/migrations/tern.conf