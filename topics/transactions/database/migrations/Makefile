-include .env
export

run-migrations:
	tern migrate -m . --config ./tern.conf

rollback-all:
	tern migrate -m . -c ./tern.conf -d 0

status:
	tern status -m . -c ./tern.conf