.PHONY: run build test test-cover test-integration docker-up docker-down tidy

run:
	go run ./cmd/api

build:
	go build -o bin/api ./cmd/api

test:
	go test ./...

test-cover:
	go test -cover ./...

test-integration:
	go test -tags=integration ./...

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

tidy:
	go mod tidy
