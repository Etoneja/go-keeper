.PHONY: test build clean lint

migrate-up:
	go run cmd/migrate/main.go up

migrate-down:
	go run cmd/migrate/main.go down

migrate-version:
	go run cmd/migrate/main.go version

test:
	@go test -v ./...

fmt:
	@go fmt ./...

build:
	@go build -o bin/server ./cmd/server/

run:
	@go run ./cmd/my-app/

lint:
	@golangci-lint run

clean:
	@rm -rf bin/

deps:
	@go mod download
	@go mod verify

all: deps test build
