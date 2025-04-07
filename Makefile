containerup:
	docker-compose up -d

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

.PHONY: sqlc containerup createdb dropdb migrateup migratedown