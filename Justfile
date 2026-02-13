# Justfile for A2A Golang Project
set shell := ["bash", "-c"]

# Run tests and linter
default: lint test

# Run all tests
test:
    go test -v ./...

# Run code linter
lint:
    golangci-lint run ./...

# Build binaries
build:
    @echo "Building Agent A (Client)..."
    go build -o bin/agent_a ./cmd/agent_a
    @echo "Building Agent B & C (Server)..."
    go build -o bin/server ./cmd/server

# Run Agent Server (B+C)
run-server:
    @echo "🚀 A2A Server (B+C) starting..."
    go run ./cmd/server/main.go

# Run Agent A (Assistant Client)
run-a:
    @echo "🚀 Agent A (Assistant Client) starting..."
    go run ./cmd/agent_a/main.go

# Clean build artifacts
clean:
    rm -rf bin/
    go clean
