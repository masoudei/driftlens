DB_URL := postgres://driftlens:driftlens@localhost:5433/driftlens?sslmode=disable

.PHONY: dev dev-db-only dev-db-stop dev-db-clean dev-db-reset dev-kind-cluster dev-kind-cluster-delete dev-kind dev-kind-db web web-build test test-cover build clean

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

# Frontend
web:
	cd web && npm install && npm run dev

web-build:
	cd web && npm run build

build:
	go build -o driftlens ./cmd/driftlens

# ArgoCD
.PHONY: dev-argocd-setup dev-argocd-clean dev-argocd-env dev-k8s

dev-argocd-setup:
	bash scripts/setup-argocd-kind.sh

dev-argocd-clean:
	-kubectl delete namespace argocd --ignore-not-found
	-kubectl delete namespace sample-app --ignore-not-found
	@echo ""
	@echo "To rebuild cluster from scratch:"
	@echo "  kind delete cluster --name driftlens"
	@echo "  kind create cluster --config kind-config.yaml --name driftlens"

dev-argocd-env:
	@echo "Required env vars for DriftLens:"
	@echo ""
	@echo "  set DRIFTLENS_ARGOCD_URL=https://localhost:8443"
	@echo "  set DRIFTLENS_ARGOCD_TOKEN= $$(argocd account generate-token)"
	@echo "  set DRIFTLENS_ARGOCD_POLL_INTERVAL=30"
	@echo ""
	@echo "Then run: go run ./cmd/driftlens"

# Run against Kind cluster with ArgoCD (no DB)
dev-k8s: build
	./driftlens

# Run against Kind cluster with ArgoCD + PostgreSQL
dev-k8s-db:
	docker compose up -d postgres
	go build -o driftlens ./cmd/driftlens
	DRIFTLENS_DATABASE_URL=$(DB_URL) ./driftlens

clean:
	rm -f driftlens
	rm -f coverage.out
