DB_URL := postgres://driftlens:driftlens@localhost:5433/driftlens?sslmode=disable

.PHONY: dev dev-db dev-db-clean dev-db-run dev-kind-cluster dev-kind-cluster-delete dev-kind dev-kind-db test test-cover build clean

# Dev mode (mock collector, in-memory store — no K8s or DB needed)
dev:
	DRIFTLENS_DEV=true go run ./cmd/driftlens

# Dev mode with PostgreSQL
dev-db-run:
	docker compose up -d postgres
	DRIFTLENS_DEV=true DRIFTLENS_DATABASE_URL=$(DB_URL) go run ./cmd/driftlens

# Stop PostgreSQL (keep data)
dev-db-stop:
	docker compose down

# Stop, wipe data, and release bound ports (use when port conflicts or auth fails)
dev-db-clean:
	docker compose down -v
	@echo ""
	@echo "If ports remain bound, restart Docker Desktop: right-click tray icon -> Restart"

# Full reset: clean + restart Docker networking on Windows
dev-db-reset: dev-db-clean

# Kind cluster management
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

# Build
build:
	go build -o driftlens ./cmd/driftlens

# Clean
clean:
	rm -f driftlens
	rm -f coverage.out
