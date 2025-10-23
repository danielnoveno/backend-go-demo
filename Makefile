build: 
	@go build -o cmd/main.go

test:
	@go test -v ./...

run: build
	@./bin/ecom

migrateion:
	@migrate create -ext sql -dir cmd/migrate/migrateion $(filter-out $@, $(MAKECMDGOALS))

migrate-up:
	@go run cmd/migrate/main.go up

migrate-down:
	@go run cmd/migrate/mian.go down