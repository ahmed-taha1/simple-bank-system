containerup:
	docker-compose up -d

containerstop:
	docker-compose stop

createdb:
	docker exec -it postgres createdb --username=taha --owner=taha simple_bank

dropdb:
	docker exec -it postgres dropdb --username=taha simple_bank

migrateup:
	migrate -path db/migration -database "postgresql://taha:root@localhost:5431/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://taha:root@localhost:5431/simple_bank?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

.PHONY: sqlc containerup createdb dropdb migrateup migratedown test containerstop server