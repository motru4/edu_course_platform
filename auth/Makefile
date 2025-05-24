.PHONY: build run test proto migrate-up migrate-down

# Build the application
build:
	go build -o bin/auth-service cmd/main.go

# Run the application
run:
	go run cmd/main.go

# Run tests
test:
	go test -v ./...

# Generate protobuf code
proto:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		internal/transport/grpc/auth.proto

# Generate swagger documentation
swag:
	swag init -g internal/app/app.go

# Database migrations
migrate-up:
	goose -dir migrations postgres "$(DB_URL)" up

migrate-down:
	goose -dir migrations postgres "$(DB_URL)" down

# Install dependencies
deps:
	go mod download
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/pressly/goose/v3/cmd/goose@latest

# Clean build artifacts
clean:
	rm -rf bin/ 