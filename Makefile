.PHONY: dev dev-kind dev-kind-cluster dev-kind-cluster-delete dev-db test test-cover build clean

# Run locally with mock collector (no cluster needed)
dev:
	DRIFTLENS_DEV=true go run ./cmd/driftlens

# Start PostgreSQL via Docker Compose
dev-db:
	docker compose up -d postgres

# Stop PostgreSQL
dev-db-stop:
	docker compose down

# Run with PostgreSQL store
dev-db-run: dev-db
	DRIFTLENS_DEV=true DRIFTLENS_DATABASE_URL="postgres://driftlens:driftlens@localhost:5432/driftlens?sslmode=disable" go run ./cmd/driftlens

# Create kind cluster (context auto-added to ~/.kube/config)
dev-kind-cluster:
	kind create cluster --config kind-config.yaml --name driftlens

# Delete kind cluster
dev-kind-cluster-delete:
	kind delete cluster --name driftlens

# Build and run, connecting to kind cluster via default kubeconfig context
dev-kind: build
	./driftlens

# Run all tests
test:
	go test ./... -v -cover

# Run tests with coverage report
test-cover:
	go test ./... -coverprofile=coverage.out -covermode=atomic
	go tool cover -func=coverage.out

# Build binary
build:
	go build -o driftlens ./cmd/driftlens

# Clean build artifacts
clean:
	rm -f driftlens
	rm -f coverage.out
