.PHONY: dev dev-kind test build clean

# Run locally with mock collector (no cluster needed)
dev:
	DRIFTLENS_DEV=true go run ./cmd/driftlens

# Create kind cluster
dev-kind-cluster:
	kind create cluster --config deploy/kind/kind-config.yaml --name driftlens

# Delete kind cluster
dev-kind-cluster-delete:
	kind delete cluster --name driftlens

# Build and run inside kind (uses in-cluster config)
dev-kind: build
	kubectl apply -f deploy/kind/
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
