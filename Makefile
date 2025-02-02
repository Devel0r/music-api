.SILENT:

.PHONY: fmt lint race test run migrate_up migrate_down migrate_status

include .env
export 

fmt:
	go fmt ./...

lint: fmt
	go vet ./...

race: lint
	go test -v ./...

test: race
	go test -v -cover ./...

run:
	go run -v cmd/music_api/main.go

migrate_up: 
	goose -dir ./migrations postgres "host=localhost port=5432 user=postgres password=postgres dbname=music_api sslmode=disable" up 

migrate_down:
	goose -dir ./migrations postgres "host=localhost port=5432 user=postgres password=postgres dbname=music_api sslmode=disable" down

migrate_status:
	goose -dir ./migrations postgres "host=localhost port=5432 user=postgres password=postgres dbname=music_api sslmode=disable" status

.DEFAULT_GOAL := run