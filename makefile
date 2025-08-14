all: build

setup:
	@go get tool
	@go install github.com/swaggo/swag/cmd/swag@latest
	@cd ui && pnpm install

build:
	@go build -o main cmd/api/main.go

run:
	@go run cmd/api/main.go

watch:
	@docker compose up backend frontend
	@docker compose down

swagger:
	@swag fmt -g ./internal/server/api/v1/v1.go
	@swag init -q -g ./internal/server/api/v1/v1.go -o ./docs/v1/

seed:
	@docker compose up --abort-on-container-exit --exit-code-from seed seed
	@docker compose down

goose:
	@docker compose down
	@docker compose up db -d
	@docker compose exec db bash -c 'until pg_isready -U postgres; do sleep 1; done'
	@read -p "Action: " action; \
	go tool goose -dir ./db/migrations postgres "user=postgres password=postgres host=localhost port=5432 dbname=events sslmode=disable" $$action
	@docker compose down db

migrate:
	@docker compose down
	@docker compose up db -d
	@docker compose exec db bash -c 'until pg_isready -U postgres; do sleep 1; done'
	@go tool goose -dir ./db/migrations postgres "user=postgres password=postgres host=localhost port=5432 dbname=events sslmode=disable" up
	@docker compose down db

create-migration:
	@read -p "Enter migration name: " name; \
	go tool goose -dir ./db/migrations create $$name sql

query:
	@go tool sqlc generate

test:
	go test ./..

dead:
	go tool deadcode ./...

.PHONY: all setup build run watch swagger seed goose migrate create-migration query test dead
