DB_URL := postgres://driftlens:driftlens@localhost:5433/driftlens?sslmode=disable

.PHONY: dev dev-db-only dev-db-stop dev-db-clean dev-db-reset dev-kind-cluster dev-kind-cluster-delete dev-kind dev-kind-db test test-cover build clean

# One command: starts PostgreSQL if not running, runs app with mock collector + DB
dev:
	docker compose up -d postgres
	DRIFTLENS_DEV=true DRIFTLENS_DATABASE_URL=$(DB_URL) go run ./cmd/driftlens

# Start PostgreSQL only (no app)
dev-db-only:
	docker compose up -d postgres

# Stop PostgreSQL (keep data)
dev-db-stop:
	docker compose down

# Stop and wipe data
dev-db-clean:
	docker compose down -v
	@echo ""
	@echo "If ports remain bound: right-click Docker Desktop tray icon -> Restart"

dev-db-reset: dev-db-clean

# Kind cluster
dev-kind-cluster:
	kind create cluster --config kind-config.yaml --name driftlens

dev-kind-cluster-delete:
	kind delete cluster --name driftlens

# Run against Kind cluster (in-memory store)
dev-kind: build
	./driftlens

# Run against Kind cluster with PostgreSQL
dev-kind-db:
	docker compose up -d postgres
	go build -o driftlens ./cmd/driftlens
	DRIFTLENS_DATABASE_URL=$(DB_URL) ./driftlens

# Tests
test:
	go test ./... -v -cover

test-cover:
	go test ./... -coverprofile=coverage.out -covermode=atomic
	go tool cover -func=coverage.out

build:
	go build -o driftlens ./cmd/driftlens

clean:
	rm -f driftlens
	rm -f coverage.out
