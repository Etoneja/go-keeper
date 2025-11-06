.PHONY: test build clean lint

migrate-up:
	go run cmd/migrate/main.go up

migrate-down:
	go run cmd/migrate/main.go down

migrate-version:
	go run cmd/migrate/main.go version

genproto:
	protoc --go_out=. --go-grpc_out=. internal/proto/api.proto --go_opt=default_api_level=API_OPAQUE

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
