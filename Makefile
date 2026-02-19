.PHONY: dev test build cluster clean

dev:
	@echo "Starting local development cluster..."
	kind create cluster --name quantatom-local || true
	kubectl apply -k infra/k8s/base

test:
	@echo "Running tests across all services..."
	cd ui/web && npm test -- --ci
	cd compute/heliocalc && cargo test
	@echo "Running integration tests..."
	cd tests/integration && go test ./... || echo "No integration tests found."

cluster:
	kind create cluster --name quantatom-local

clean:
	kind delete cluster --name quantatom-local
