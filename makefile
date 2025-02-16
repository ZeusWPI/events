.PHONY: build watch db sqlc create-migration goose migrate test

build:
	@go build -o main cmd/events/main.go

watch: 
	@docker compose up

sqlc:
	sqlc generate

create-migration:
	@read -p "Enter migration name: " name; \
		goose -dir ./db/migrations create $$name sql 

goose:
	@read -p "Action: " action; \
	goose -dir ./db/migrations postgres "user=postgres password=postgres dbname=events host=localhost sslmode=disable" $$action

migrate:
	@goose -dir ./db/migrations postgres "user=postgres password=postgres dbname=events host=localhost sslmode=disable" up

test:
	go test ./...
