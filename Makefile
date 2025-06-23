# Makefile for hep-sidekick

BINARY_NAME=hep-sidekick
CMD_PATH=./cmd/hep-sidekick

.PHONY: all build clean run help

all: help

# Build the binary
build:
	@echo "Building $(BINARY_NAME)..."
	@go build -o $(BINARY_NAME) $(CMD_PATH)

# Clean the binary
clean:
	@echo "Cleaning..."
	@go clean
	@rm -f $(BINARY_NAME)

# Run the binary
run: build
	@echo "Running $(BINARY_NAME)..."
	@./$(BINARY_NAME)

# Help
help:
	@echo "Usage: make [target]"
	@echo "targets:"
	@echo "  build    - build the binary"
	@echo "  clean    - clean the binary"
	@echo "  run      - run the binary"
	@echo "  help     - show this help" 