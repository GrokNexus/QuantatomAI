.PHONY: build test run dev

build:
	docker build -t $(shell basename $(CURDIR)) .

test:
	go test ./... || npm test || cargo test

run:
	go run src/main.go || npm start || cargo run

dev:
	@echo "Starting $(shell basename $(CURDIR)) in dev mode..."
	# Placeholder for tilt/skaffold or local run
	$(MAKE) run
