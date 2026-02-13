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
    @echo "Building Agent A..."
    go build -o bin/agent_a ./cmd/agent_a
    @echo "Building Agent B..."
    go build -o bin/agent_b ./cmd/agent_b

# Run Agent B (Finance Server)
run-b:
    @echo "🚀 Agent B (Finance Server) starting..."
    go run ./cmd/agent_b/main.go

# Run Agent A (Assistant Client)
run-a:
    @echo "🚀 Agent A (Assistant Client) starting..."
    go run ./cmd/agent_a/main.go

# Clean build artifacts
clean:
    rm -rf bin/
    go clean
