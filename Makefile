-include .env
export

run-migrations:
	tern migrate -m ./database/migrations/ --config ./database/migrations/tern.conf

rollback-all:
	tern migrate -m ./database/migrations -c ./database/migrations/tern.conf -d 0

status:
	tern status -m ./database/migrations -c ./database/migrations/tern.conf