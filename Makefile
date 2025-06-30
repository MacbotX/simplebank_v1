postgres:
	docker run --name simplebank_postgres -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:17.4-alpine3.21

createImage:
 	docker run --name simplebank --network bank-network -p 8080:8080 -e GIN_MODE=release -e DB_SOURCE="postgresql://root:secret@simplebank_postgres:5432/simple_bank?sslmode=disable" simplebank:latest

createdb:
	docker exec -it simplebank_postgres createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it simplebank dropdb simple_bank

migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up

migrateup1:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up 1

migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down

migratedown1:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down 1

# to generate sqlc migrations
sqlc:
	sqlc generate

# to test go file and funcs 
test:
	go test -v -cover ./...

# to tidy the code base by installing dependencies
tidy:
	go mod tidy

server:
	go run main.go

mock:
	mockgen -destination db/mock/store.go  github.com/MacbotX/simplebank_v1/db/sqlc Store

.PHONY: createdb dropdb postgres migrateup migratedown sqlc tidy test server mock migrateup1 migratedown1