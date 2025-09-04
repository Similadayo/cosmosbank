# Project variables
BINARY_NAME = cosmosbankd
GOJSON = 1

.PHONY: build run test clean

## Build the blockchain binary
build:
	@echo "Building $(BINARY_NAME)..."
	GOJSON=$(GOJSON) go build -o build/$(BINARY_NAME) ./cmd/cosmosbankd

## Run the blockchain (single-node testnet)
run: build
	@echo "Starting $(BINARY_NAME)..."
	./build/$(BINARY_NAME) start

## Run all tests with GOJSON=1 to avoid Sonic issues
test:
	@echo "Running tests..."
	GOJSON=$(GOJSON) go test ./... -v

## Clean build artifacts
clean:
	@echo "Cleaning up..."
	rm -rf build
