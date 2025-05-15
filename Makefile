postgres:
	docker run --name simplebank -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:17.4-alpine3.21

createdb:
	docker exec -it simplebank createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it simplebank dropdb simple_bank

migrateup:
	# migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up
	until migrate -path db/migration -database "postgresql://root:secret@postgres:5432/simple_bank?sslmode=disable" -verbose up; do echo "Waiting for DB..."; sleep 2; done

migratedown:
	# migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down
	until migrate -path db/migration -database "postgresql://root:secret@postgres:5432/simple_bank?sslmode=disable" -verbose down; do echo "Waiting for DB..."; sleep 2; done

# to generate sqlc migrations
sqlc:
	sqlc generate

# to test go file and funcs 
test:
	go test -v -cover ./...

# to tidy the code base by installing dependencies
tidy:
	go mod tidy

.PHONY: createdb dropdb postgres migrateup migratedown sqlc tidy test