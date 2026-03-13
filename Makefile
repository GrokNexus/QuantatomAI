.PHONY: dev test build cluster clean clean-artifacts install

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

# Install all dependencies from scratch (after a clean or fresh clone)
install:
	@echo "Installing UI dependencies..."
	cd ui/web && npm install
	@echo "Fetching Go modules..."
	cd services/grid-service && go mod download
	cd services/cortex-service && go mod download || true
	cd services/modeling-service && go mod download || true
	cd services/planning-service && go mod download || true
	@echo "Done. Run 'make dev' to start."

# Delete generated build artifacts (node_modules, .next, Rust target/, Go bin/)
# Safe to run any time — re-run 'make install' and 'make build' to regenerate.
clean-artifacts:
	@echo "Removing UI build artifacts..."
	rm -rf ui/web/node_modules ui/web/.next ui/web/dist ui/web/out
	@echo "Removing Rust build artifacts..."
	rm -rf services/atom-engine/target compute/heliocalc/target
	@echo "Removing Go build outputs..."
	rm -rf services/grid-service/bin services/cortex-service/bin
	@echo "Cleaned. Disk space recovered."

# Full clean: artifacts + kind cluster
clean: clean-artifacts
	kind delete cluster --name quantatom-local || true
